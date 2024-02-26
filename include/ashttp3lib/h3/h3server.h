// Copytight (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included below.
//
// Copyright (C) 2018-2019, Cloudflare, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
// THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
// EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
// PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

#include <inttypes.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string>

#include <errno.h>
#include <fcntl.h>

#include <netdb.h>
#include <sys/socket.h>
#include <sys/types.h>

#include <ev.h>
#include <uthash.h>

#include <quiche.h>

#include <h3request.h>
#include <h3response.h>
#include <functional>
#include <unordered_map>

namespace ashttp3lib {
const int LOCAL_CONN_ID_LEN = 16;
const int MAX_DATAGRAM_SIZE = 1350;
const int MAX_TOKEN_LEN = sizeof("ashttp3lib") - 1 +
                          sizeof(struct sockaddr_storage) +
                          QUICHE_MAX_CONN_ID_LEN;

struct connections {
  int sock;

  struct sockaddr* local_addr;
  socklen_t local_addr_len;

  struct conn_io* h;
};

struct conn_io {
  ev_timer timer;

  int sock;

  uint8_t cid[LOCAL_CONN_ID_LEN];

  quiche_conn* conn;
  quiche_h3_conn* http3;

  struct sockaddr_storage peer_addr;
  socklen_t peer_addr_len;

  UT_hash_handle hh;
};

static quiche_config* config = NULL;

static quiche_h3_config* http3_config = NULL;

static struct connections* conns = NULL;

typedef std::function<std::string(ashttp3lib::H3request&)> CallbackFunction;

class Http3Server {
 private:
  struct addrinfo* local;
  int sock;
  std::string host;
  std::string port;
  std::unordered_map<
      std::string,
      std::unordered_map<std::string,
                         std::function<std::string(ashttp3lib::H3request&)>>>
      routes_;

 public:
  Http3Server(std::string, std::string, bool);
  void run();

 private:
  void timeout_cb(EV_P_ ev_timer* w, int revents);
  void recv_cb(EV_P_ ev_io* w, int revents);
  static int for_each_header(uint8_t* name, size_t name_len, uint8_t* value,
                             size_t value_len, void* argp);
  static struct conn_io* create_conn(uint8_t* scid, size_t scid_len,
                                     uint8_t* odcid, size_t odcid_len,
                                     struct sockaddr* local_addr,
                                     socklen_t local_addr_len,
                                     struct sockaddr_storage* peer_addr,
                                     socklen_t peer_addr_len);
  uint8_t* gen_cid(uint8_t* cid, size_t cid_len);
  bool validate_token(const uint8_t* token, size_t token_len,
                      struct sockaddr_storage* addr, socklen_t addr_len,
                      uint8_t* odcid, size_t* odcid_len);
  void mint_token(const uint8_t* dcid, size_t dcid_len,
                  struct sockaddr_storage* addr, socklen_t addr_len,
                  uint8_t* token, size_t* token_len);
  void flush_egress(struct ev_loop* loop, struct conn_io* conn_io);
  static void debug_log(const char* line, void* argp);
  void add_route(std::string method, std::string path,
                 CallbackFunction request);

  // helper static functions, wrapper around member methods
  static void static_recv_cb(struct ev_loop* loop, ev_io* w, int revents) {
    Http3Server* server = static_cast<Http3Server*>(w->data);
    server->recv_cb(loop, w, revents);
  }

  static void static_timeout_cb(struct ev_loop* loop, ev_timer* w,
                                int revents) {
    Http3Server* server = static_cast<Http3Server*>(w->data);
    server->timeout_cb(loop, w, revents);
  }
};

Http3Server::Http3Server(std::string host_in, std::string port_in,
                         bool enable_debug) {
  this->host = std::move(host_in);
  this->port = std::move(port_in);
  const struct addrinfo hints = {.ai_family = PF_UNSPEC,
                                 .ai_socktype = SOCK_DGRAM,
                                 .ai_protocol = IPPROTO_UDP};

  if (enable_debug)
    quiche_enable_debug_logging(debug_log, NULL);

  if (getaddrinfo(host.c_str(), port.c_str(), &hints, &local) != 0) {
    perror("failed to resolve host");
    exit(-1);
  }

  this->sock = socket(local->ai_family, SOCK_DGRAM, 0);
  if (sock < 0) {
    perror("failed to create socket");
    exit(-1);
  }

  if (fcntl(sock, F_SETFL, O_NONBLOCK) != 0) {
    perror("failed to make socket non-blocking");
    exit(-1);
  }

  if (bind(sock, local->ai_addr, local->ai_addrlen) < 0) {
    perror("failed to connect socket");
    exit(-1);
  }

  config = quiche_config_new(QUICHE_PROTOCOL_VERSION);
  if (config == NULL) {
    fprintf(stderr, "failed to create config\n");
    exit(-1);
  }

  quiche_config_load_cert_chain_from_pem_file(config, "./cert.crt");
  quiche_config_load_priv_key_from_pem_file(config, "./cert.key");

  quiche_config_set_application_protos(
      config, (uint8_t*)QUICHE_H3_APPLICATION_PROTOCOL,
      sizeof(QUICHE_H3_APPLICATION_PROTOCOL) - 1);

  quiche_config_set_max_idle_timeout(config, 5000);
  quiche_config_set_max_recv_udp_payload_size(config, MAX_DATAGRAM_SIZE);
  quiche_config_set_max_send_udp_payload_size(config, MAX_DATAGRAM_SIZE);
  quiche_config_set_initial_max_data(config, 10000000);
  quiche_config_set_initial_max_stream_data_bidi_local(config, 1000000);
  quiche_config_set_initial_max_stream_data_bidi_remote(config, 1000000);
  quiche_config_set_initial_max_stream_data_uni(config, 1000000);
  quiche_config_set_initial_max_streams_bidi(config, 100);
  quiche_config_set_initial_max_streams_uni(config, 100);
  quiche_config_set_disable_active_migration(config, true);
  quiche_config_set_cc_algorithm(config, QUICHE_CC_RENO);

  http3_config = quiche_h3_config_new();
  if (http3_config == NULL) {
    fprintf(stderr, "failed to create HTTP/3 config\n");
    exit(-1);
  }
}

void Http3Server::run() {
  struct connections c;
  c.sock = this->sock;
  c.h = NULL;
  c.local_addr = this->local->ai_addr;
  c.local_addr_len = this->local->ai_addrlen;

  conns = &c;

  ev_io watcher;

  struct ev_loop* loop = ev_default_loop(0);

  ev_io_init(&watcher, static_recv_cb, this->sock, EV_READ);
  ev_io_start(loop, &watcher);
  watcher.data = &c;

  ev_loop(loop, 0);

  freeaddrinfo(this->local);

  quiche_h3_config_free(http3_config);

  quiche_config_free(config);
}

void Http3Server::timeout_cb(EV_P_ ev_timer* w, int revents) {
  struct conn_io* conn_io = (struct conn_io*)(w->data);
  quiche_conn_on_timeout(conn_io->conn);

  fprintf(stderr, "timeout\n");

  flush_egress(loop, conn_io);

  if (quiche_conn_is_closed(conn_io->conn)) {
    quiche_stats stats;
    quiche_path_stats path_stats;

    quiche_conn_stats(conn_io->conn, &stats);
    quiche_conn_path_stats(conn_io->conn, 0, &path_stats);

    fprintf(stderr,
            "connection closed, recv=%zu sent=%zu lost=%zu rtt=%" PRIu64
            "ns cwnd=%zu\n",
            stats.recv, stats.sent, stats.lost, path_stats.rtt,
            path_stats.cwnd);

    HASH_DELETE(hh, conns->h, conn_io);

    ev_timer_stop(loop, &conn_io->timer);
    quiche_conn_free(conn_io->conn);
    free(conn_io);

    return;
  }
}

void Http3Server::recv_cb(EV_P_ ev_io* w, int revents) {
  struct conn_io *tmp, *conn_io = NULL;

  static uint8_t buf[65535];
  static uint8_t out[MAX_DATAGRAM_SIZE];

  while (1) {
    struct sockaddr_storage peer_addr;
    socklen_t peer_addr_len = sizeof(peer_addr);
    memset(&peer_addr, 0, peer_addr_len);

    ssize_t read = recvfrom(conns->sock, buf, sizeof(buf), 0,
                            (struct sockaddr*)&peer_addr, &peer_addr_len);

    if (read < 0) {
      if ((errno == EWOULDBLOCK) || (errno == EAGAIN)) {
        fprintf(stderr, "recv would block\n");
        break;
      }

      perror("failed to read");
      return;
    }

    uint8_t type;
    uint32_t version;

    uint8_t scid[QUICHE_MAX_CONN_ID_LEN];
    size_t scid_len = sizeof(scid);

    uint8_t dcid[QUICHE_MAX_CONN_ID_LEN];
    size_t dcid_len = sizeof(dcid);

    uint8_t odcid[QUICHE_MAX_CONN_ID_LEN];
    size_t odcid_len = sizeof(odcid);

    uint8_t token[MAX_TOKEN_LEN];
    size_t token_len = sizeof(token);

    int rc =
        quiche_header_info(buf, read, LOCAL_CONN_ID_LEN, &version, &type, scid,
                           &scid_len, dcid, &dcid_len, token, &token_len);
    if (rc < 0) {
      fprintf(stderr, "failed to parse header: %d\n", rc);
      return;
    }

    HASH_FIND(hh, conns->h, dcid, dcid_len, conn_io);

    if (conn_io == NULL) {
      if (!quiche_version_is_supported(version)) {
        fprintf(stderr, "version negotiation\n");

        ssize_t written = quiche_negotiate_version(scid, scid_len, dcid,
                                                   dcid_len, out, sizeof(out));

        if (written < 0) {
          fprintf(stderr, "failed to create vneg packet: %zd\n", written);
          continue;
        }

        ssize_t sent = sendto(conns->sock, out, written, 0,
                              (struct sockaddr*)&peer_addr, peer_addr_len);
        if (sent != written) {
          perror("failed to send");
          continue;
        }

        fprintf(stderr, "sent %zd bytes\n", sent);
        continue;
      }

      if (token_len == 0) {
        fprintf(stderr, "stateless retry\n");

        mint_token(dcid, dcid_len, &peer_addr, peer_addr_len, token,
                   &token_len);

        uint8_t new_cid[LOCAL_CONN_ID_LEN];

        if (gen_cid(new_cid, LOCAL_CONN_ID_LEN) == NULL) {
          continue;
        }

        ssize_t written = quiche_retry(scid, scid_len, dcid, dcid_len, new_cid,
                                       LOCAL_CONN_ID_LEN, token, token_len,
                                       version, out, sizeof(out));

        if (written < 0) {
          fprintf(stderr, "failed to create retry packet: %zd\n", written);
          continue;
        }

        ssize_t sent = sendto(conns->sock, out, written, 0,
                              (struct sockaddr*)&peer_addr, peer_addr_len);
        if (sent != written) {
          perror("failed to send");
          continue;
        }

        fprintf(stderr, "sent %zd bytes\n", sent);
        continue;
      }

      if (!validate_token(token, token_len, &peer_addr, peer_addr_len, odcid,
                          &odcid_len)) {
        fprintf(stderr, "invalid address validation token\n");
        continue;
      }

      conn_io = create_conn(dcid, dcid_len, odcid, odcid_len, conns->local_addr,
                            conns->local_addr_len, &peer_addr, peer_addr_len);

      if (conn_io == NULL) {
        continue;
      }
    }

    quiche_recv_info recv_info = {
        (struct sockaddr*)&peer_addr,
        peer_addr_len,

        conns->local_addr,
        conns->local_addr_len,
    };

    ssize_t done = quiche_conn_recv(conn_io->conn, buf, read, &recv_info);

    if (done < 0) {
      fprintf(stderr, "failed to process packet: %zd\n", done);
      continue;
    }

    fprintf(stderr, "recv %zd bytes\n", done);

    if (quiche_conn_is_established(conn_io->conn)) {
      quiche_h3_event* ev;

      if (conn_io->http3 == NULL) {
        conn_io->http3 =
            quiche_h3_conn_new_with_transport(conn_io->conn, http3_config);
        if (conn_io->http3 == NULL) {
          fprintf(stderr, "failed to create HTTP/3 connection\n");
          continue;
        }
      }

      while (1) {
        // each connection can be idenfied with this, therefore process anything
        // related to the request with this string. maybe construct an asio
        // queue of response object and send them
        int64_t s = quiche_h3_conn_poll(conn_io->http3, conn_io->conn, &ev);

        if (s < 0) {
          break;
        }

        // define a new H3Request Object, and pass it down
        // to get the headers, and body of the request
        ashttp3lib::H3request request;
        ashttp3lib::H3response response;
        switch (quiche_h3_event_type(ev)) {
          case QUICHE_H3_EVENT_HEADERS: {
            // an event loop handles parsing of headers -> asynchronous processing
            int rc =
                quiche_h3_event_for_each_header(ev, for_each_header, &request);

            if (rc != 0) {
              response.set_status("422");
              response.set_body(
                  "Unprocessable Entity.");
              fprintf(stderr, "failed to process headers\n");
            }

            break;
          }

          case QUICHE_H3_EVENT_DATA: {
            // TODO: parse the body also of the incoming request
            fprintf(stderr, "got HTTP data\n");
            break;
          }

          case QUICHE_H3_EVENT_FINISHED:
            if (this->routes_.find(request.get_path()) == this->routes_.end()) {
              response.set_status("404");
              response.set_body("Not Found");
            } else if (this->routes_.at(request.get_path())
                           .find(request.get_method()) ==
                       this->routes_.at(request.get_path()).end()) {
              response.set_status("405");
              response.set_body("Method Not Allowed");
            } else if (response.is_ok()) {
              // process the bound function only if there's no error
              response.set_status("200");
              response.set_body(routes_.at(request.get_path())
                                    .at(request.get_method())(request));
            }

            response.add_headers("server", "ashttp3lib");
            response.add_headers("content-length",
                                 std::to_string(response.get_content_len()));

            quiche_h3_send_response(conn_io->http3, conn_io->conn, s,
                                    response.get_headers(),
                                    response.get_header_len(), false);

            quiche_h3_send_body(conn_io->http3, conn_io->conn, s,
                                (uint8_t*)response.serialize_response().c_str(),
                                sizeof(response), true);
            break;

          case QUICHE_H3_EVENT_RESET:
            break;

          case QUICHE_H3_EVENT_PRIORITY_UPDATE:
            break;

          case QUICHE_H3_EVENT_GOAWAY: {
            fprintf(stderr, "got GOAWAY\n");
            break;
          }
        }

        // deallocate the request object constructed
        // delete request;
        quiche_h3_event_free(ev);
      }
    }
  }

  HASH_ITER(hh, conns->h, conn_io, tmp) {
    flush_egress(loop, conn_io);

    if (quiche_conn_is_closed(conn_io->conn)) {
      quiche_stats stats;
      quiche_path_stats path_stats;

      quiche_conn_stats(conn_io->conn, &stats);
      quiche_conn_path_stats(conn_io->conn, 0, &path_stats);

      fprintf(stderr,
              "connection closed, recv=%zu sent=%zu lost=%zu rtt=%" PRIu64
              "ns cwnd=%zu\n",
              stats.recv, stats.sent, stats.lost, path_stats.rtt,
              path_stats.cwnd);

      HASH_DELETE(hh, conns->h, conn_io);

      ev_timer_stop(loop, &conn_io->timer);

      quiche_conn_free(conn_io->conn);
      free(conn_io);
    }
  }
}

int Http3Server::for_each_header(uint8_t* name, size_t name_len, uint8_t* value,
                                 size_t value_len, void* argp) {
  // parse the headers and add them to request object
  ashttp3lib::H3request* request = (ashttp3lib::H3request*)argp;
  request->add_headers(
      std::string(reinterpret_cast<const char*>(name), name_len),
      std::string(reinterpret_cast<const char*>(value), value_len));

  fprintf(stderr, "got HTTP header: %.*s=%.*s\n", (int)name_len, name,
          (int)value_len, value);

  return 0;
}

struct conn_io* Http3Server::create_conn(uint8_t* scid, size_t scid_len,
                                         uint8_t* odcid, size_t odcid_len,
                                         struct sockaddr* local_addr,
                                         socklen_t local_addr_len,
                                         struct sockaddr_storage* peer_addr,
                                         socklen_t peer_addr_len) {
  struct conn_io* conn_io = (struct conn_io*)calloc(1, sizeof(*conn_io));
  if (conn_io == NULL) {
    fprintf(stderr, "failed to allocate connection IO\n");
    return NULL;
  }

  if (scid_len != LOCAL_CONN_ID_LEN) {
    fprintf(stderr, "failed, scid length too short\n");
  }

  memcpy(conn_io->cid, scid, LOCAL_CONN_ID_LEN);

  quiche_conn* conn = quiche_accept(
      conn_io->cid, LOCAL_CONN_ID_LEN, odcid, odcid_len, local_addr,
      local_addr_len, (struct sockaddr*)peer_addr, peer_addr_len, config);

  if (conn == NULL) {
    fprintf(stderr, "failed to create connection\n");
    return NULL;
  }

  conn_io->sock = conns->sock;
  conn_io->conn = conn;

  memcpy(&conn_io->peer_addr, peer_addr, peer_addr_len);
  conn_io->peer_addr_len = peer_addr_len;

  ev_init(&conn_io->timer, static_timeout_cb);
  conn_io->timer.data = conn_io;

  HASH_ADD(hh, conns->h, cid, LOCAL_CONN_ID_LEN, conn_io);

  fprintf(stderr, "new connection\n");

  return conn_io;
}

uint8_t* Http3Server::gen_cid(uint8_t* cid, size_t cid_len) {
  int rng = open("/dev/urandom", O_RDONLY);
  if (rng < 0) {
    perror("failed to open /dev/urandom");
    return NULL;
  }

  ssize_t rand_len = read(rng, cid, cid_len);
  if (rand_len < 0) {
    perror("failed to create connection ID");
    return NULL;
  }

  return cid;
}

bool Http3Server::validate_token(const uint8_t* token, size_t token_len,
                                 struct sockaddr_storage* addr,
                                 socklen_t addr_len, uint8_t* odcid,
                                 size_t* odcid_len) {
  if ((token_len < sizeof("quiche") - 1) ||
      memcmp(token, "quiche", sizeof("quiche") - 1)) {
    return false;
  }

  token += sizeof("quiche") - 1;
  token_len -= sizeof("quiche") - 1;

  if ((token_len < addr_len) || memcmp(token, addr, addr_len)) {
    return false;
  }

  token += addr_len;
  token_len -= addr_len;

  if (*odcid_len < token_len) {
    return false;
  }

  memcpy(odcid, token, token_len);
  *odcid_len = token_len;

  return true;
}

void Http3Server::mint_token(const uint8_t* dcid, size_t dcid_len,
                             struct sockaddr_storage* addr, socklen_t addr_len,
                             uint8_t* token, size_t* token_len) {
  memcpy(token, "quiche", sizeof("quiche") - 1);
  memcpy(token + sizeof("quiche") - 1, addr, addr_len);
  memcpy(token + sizeof("quiche") - 1 + addr_len, dcid, dcid_len);

  *token_len = sizeof("quiche") - 1 + addr_len + dcid_len;
}

void Http3Server::flush_egress(struct ev_loop* loop, struct conn_io* conn_io) {
  static uint8_t out[MAX_DATAGRAM_SIZE];

  quiche_send_info send_info;

  while (1) {
    ssize_t written =
        quiche_conn_send(conn_io->conn, out, sizeof(out), &send_info);

    if (written == QUICHE_ERR_DONE) {
      fprintf(stderr, "done writing\n");
      break;
    }

    if (written < 0) {
      fprintf(stderr, "failed to create packet: %zd\n", written);
      return;
    }

    ssize_t sent =
        sendto(conn_io->sock, out, written, 0,
               (struct sockaddr*)&conn_io->peer_addr, conn_io->peer_addr_len);
    if (sent != written) {
      perror("failed to send");
      return;
    }

    fprintf(stderr, "sent %zd bytes\n", sent);
  }

  double t = quiche_conn_timeout_as_nanos(conn_io->conn) / 1e9f;
  conn_io->timer.repeat = t;
  ev_timer_again(loop, &conn_io->timer);
}

void Http3Server::debug_log(const char* line, void* argp) {
  fprintf(stderr, "%s\n", line);
}

void Http3Server::add_route(std::string method, std::string path,
                            CallbackFunction bind_func) {
  routes_[path][method] = bind_func;
}

}  // namespace ashttp3lib
