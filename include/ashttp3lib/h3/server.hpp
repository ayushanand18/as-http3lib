#include <cstdlib>
#include <cstddef>
#include <cstdint>
#include <cstdio>
#include <memory>

#include "quiche/quic/core/crypto/crypto_server_options.h"
#include "quiche/quic/core/http/quic_simple_server.h"
#include "quiche/quic/test_tools/quic_test_server.h"
#include "quiche/quic/test_tools/quic_test_server_utils.h"
#include "quiche/quic/tools/quic_simple_crypto_server_stream_helper.h"
#include "quiche/quic/tools/quic_simple_server_session_helper.h"

namespace {

class MyHttpServerSession : public quic::QuicSimpleServerSession {
 public:
  MyHttpServerSession(const quic::QuicConfig& config,
                     quic::QuicConnection* connection,
                     quic::QuicCryptoServerStreamBase::Helper* helper,
                     const quic::QuicCryptoServerConfig* crypto_config)
      : QuicSimpleServerSession(config, connection, helper, crypto_config) {}

  ~MyHttpServerSession() override = default;

  void OnRequest(quic::QuicStreamId stream_id,
                 const quic::SpdyHeaderBlock& headers) override {
    // Handle the incoming request and send a response.
    std::string path = headers.GetHeader(":path");
    quic::QuicHttpResponseFrame response(stream_id, true);
    std::string body = "Hello, " + path + "!";
    response.set_body(body);
    SendHttpFrame(response);
  }
};

class QuicHttpServer : public quic::QuicSimpleServer {
 public:
  QuicHttpServer(quic::QuicEpollServer* epoll_server,
                 const quic::QuicConfig& quic_config,
                 const quic::ParsedQuicVersionVector& supported_versions)
      : QuicSimpleServer(epoll_server, quic_config, supported_versions) {}

  ~QuicHttpServer() override = default;

 private:
  std::unique_ptr<quic::QuicSession>
  CreateQuicSession(quic::QuicConnection* connection,
                    quic::QuicCryptoServerStreamBase::Helper* helper,
                    const quic::QuicCryptoServerConfig* crypto_config,
                    const quic::QuicConfig& config) override {
    return std::make_unique<MyHttpServerSession>(
        config, connection, helper, crypto_config);
  }
};

}  // namespace

int main() {
  quic::QuicEnableDebugLogging(quic::kQuicLogDebug, quic::kDefaultDebugCategory);

  quic::QuicEpollServer epoll_server;

  quic::QuicConfig quic_config;
  quic::QuicCryptoServerConfig crypto_config(
      quic::QuicRandom::GetInstance()->RandUint64(),
      quic::QuicRandom::GetInstance(),
      std::make_unique<quic::ProofSourceChromium>(),
      quic::KeyExchangeSource::Default());
  quic::ParsedQuicVersionVector supported_versions = {
      quic::ParsedQuicVersion(quic::PROTOCOL_QUIC_H3, quic::QUIC_VERSION_99)};

  QuicHttpServer server(&epoll_server, quic_config, supported_versions);

  quic::QuicSocketAddress bind_address("127.0.0.1", 4433);
  if (!server.Listen(bind_address)) {
    fprintf(stderr, "Failed to listen on %s\n", bind_address.ToString().c_str());
    return EXIT_FAILURE;
  }

  fprintf(stdout, "Server listening on %s\n", bind_address.ToString().c_str());

  epoll_server.Run();

  return EXIT_SUCCESS;
}
