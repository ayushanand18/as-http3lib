/*
  ashttp3lib/h1/server.hpp - A C++ HTTP/1.1 Library using Boost.Asio Server Class
  
  Copyright (c) 2024, Ayush Anand
  
  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:
  
  The above copyright notice and this permission notice shall be included in all
  copies or substantial portions of the Software.
  
  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  SOFTWARE.
*/

#include <boost/asio.hpp>
#include <iostream>
#include <string>
#include <ashttp3lib/h1/logging.hpp>
#include <ashttp3lib/h1/request.hpp>

namespace ashttp3lib::h1 {

//! \brief HTTPServer Class. A simple HTTP/1.1 Server.
class HTTPServer {
  ashttp3lib::logging::Logger* logger;

 public:
  //! \brief Constructor for HTTPServer class.
  //! \param port_num. [int] Port number for the server to bind.
  HTTPServer(int port_num)
      : acceptor(io, boost::asio::ip::tcp::endpoint(boost::asio::ip::tcp::v4(),
                                                    port_num)),
        socket(io) {
    this->logger = new ashttp3lib::logging::Logger();
  }

  //! \brief Run Server. Make a blocking run of the server instance.
  //! Listen to requests and process them.
  void run() {
    acceptRequest();
    io.run();
  }

  //! \brief Register a GET handler.
  //! \param path. [std::string] Path on the server to register GET handler.
  //! \param bind_func. [std::function<std::string(Request&)>] Callback function to handle GET requests.
  void get(std::string path,
           std::function<std::string(ashttp3lib::h1::Request&)> bind_func) {
    routes_[path]["GET"] = bind_func;
  }

  //! \brief Register a POST handler.
  //! \param path. [std::string] Path on the server to register POST handler.
  //! \param bind_func. [std::function<std::string(Request&)>] Callback function to handle POST requests.
  void post(std::string path,
            std::function<std::string(ashttp3lib::h1::Request&)> bind_func) {
    routes_[path]["POST"] = bind_func;
  }

 private:
  //! \brief Asynchronously accept incoming requests.
  void acceptRequest() {
    acceptor.async_accept(socket, [this](const boost::system::error_code& ec) {
      if (!ec) {
        handleRequest();
      }
      acceptRequest();
    });
  }

  //! \brief Asynchronously handle incoming requests.
  void handleRequest() {
    boost::asio::async_read_until(socket, request, "\r\n\r\n",
                                  [this](const boost::system::error_code& ec,
                                         std::size_t bytes_transferred) {
                                    if (!ec) {
                                      processRequest();
                                    }
                                  });
  }

  //! \brief Process the incoming request.
  void processRequest() {
    std::string result(boost::asio::buffer_cast<const char*>(request.data()),
                       boost::asio::buffer_size(request.data()));
    auto request_packet = ashttp3lib::h1::Request(result);

    mapRequestWithResponse(request_packet);
  }

  //! \brief Map the incoming request to an appropriate response handler.
  void mapRequestWithResponse(Request& request_packet) {
    if (routes_.find(request_packet.path) == routes_.end()) {
      this->logger->info(request_packet.method + " " + request_packet.path +
                         " 404 Not Found");
      sendResponse("404 Not Found", "The resource was not found on server.");
    } else if (routes_[request_packet.path].find(request_packet.method) ==
               routes_[request_packet.path].end()) {
      this->logger->info(request_packet.method + " " + request_packet.path +
                         " 405 Method Not Allowed");
      sendResponse("405 Method Not Allowed", "The used method is not allowed.");
    } else {
      auto response =
          routes_[request_packet.path][request_packet.method](request_packet);
      this->logger->info(request_packet.method + " " + request_packet.path +
                         " 200 OK");
      sendResponse("200 OK", response);
    }
  }

  //! \brief Send a response based on the processed information.
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

  //! Boost-based data members for IO Operations.
  boost::asio::io_service io;
  boost::asio::ip::tcp::acceptor acceptor;
  boost::asio::ip::tcp::socket socket;
  boost::asio::streambuf request;
  boost::asio::streambuf response;

  //! A Routes/Methods Handler
  std::unordered_map<
      std::string,
      std::unordered_map<std::string,
                         std::function<std::string(ashttp3lib::h1::Request&)>>>
      routes_;
};

}  // namespace ashttp3lib::h1

// ashttp3lib/h1/server.hpp