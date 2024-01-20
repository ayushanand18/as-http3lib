#include <boost/asio.hpp>
#include <iostream>
#include <string>
#include "logging.hpp"
#include "request.hpp"


namespace ashttp3lib::h1 {

class HTTPServer {
  ashttp3lib::logging::Logger* logger;

 public:
  HTTPServer(int port_num)
      : acceptor(io, boost::asio::ip::tcp::endpoint(boost::asio::ip::tcp::v4(),
                                                    port_num)),
        socket(io) {
    this->logger = new ashttp3lib::logging::Logger();
  }

  void run() {
    acceptRequest();
    io.run();
  }

  void get(std::string path,
           std::function<std::string(ashttp3lib::h1::Request&)> bind_func) {
    routes_[path]["GET"] = bind_func;
  }

  void post(std::string path,
            std::function<std::string(ashttp3lib::h1::Request&)> bind_func) {
    routes_[path]["POST"] = bind_func;
  }

 private:
  void acceptRequest() {
    acceptor.async_accept(socket, [this](const boost::system::error_code& ec) {
      if (!ec) {
        handleRequest();
      }
      acceptRequest();
    });
  }

  void handleRequest() {
    boost::asio::async_read_until(socket, request, "\r\n\r\n",
                                  [this](const boost::system::error_code& ec,
                                         std::size_t bytes_transferred) {
                                    if (!ec) {
                                      processRequest();
                                    }
                                  });
  }

  void processRequest() {
    std::istream request_stream(&request);
    auto request_packet = ashttp3lib::h1::Request(request_stream);

    mapRequestWithResponse(request_packet);
  }

  void mapRequestWithResponse(Request& request_packet) {
    if (routes_.find(request_packet.path) == routes_.end()) {
      this->logger->info(request_packet.method + " " + request_packet.path + " 404 Not Found");
      sendResponse("404 Not Found", "The resource was not found on server.");
    } else if (routes_[request_packet.path].find(request_packet.method) ==
               routes_[request_packet.path].end()) {
      this->logger->info(request_packet.method + " " + request_packet.path + " 405 Method Not Allowed");
      sendResponse("405 Method Not Allowed", "The used method is not allowed.");
    } else {
      auto response =
          routes_[request_packet.path][request_packet.method](request_packet);
      this->logger->info(request_packet.method + " " + request_packet.path + " 200 OK");
      sendResponse("200 OK", response);
    }
  }

  void sendResponse(const std::string& status, const std::string& content) {
    std::ostream response_stream(&response);
    response_stream << "HTTP/1.1 " << status << "\r\n";
    response_stream << "Content-Length: " << content.length() << "\r\n";
    response_stream << "Content-Type: text/plain\r\n\r\n";
    response_stream << content;

    boost::asio::async_write(
        socket, response,
        [this](const boost::system::error_code& ec,
               std::size_t bytes_transferred) { socket.close(); });
  }

  boost::asio::io_service io;
  boost::asio::ip::tcp::acceptor acceptor;
  boost::asio::ip::tcp::socket socket;
  boost::asio::streambuf request;
  boost::asio::streambuf response;

  std::unordered_map<
      std::string,
      std::unordered_map<std::string,
                         std::function<std::string(ashttp3lib::h1::Request&)>>>
      routes_;
};

}  // namespace ashttp3lib::h1