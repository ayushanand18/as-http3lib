/*
  ashttp3lib/h1/request.hpp - A C++ HTTP/1.1 Library Request Class
  
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

#include <iostream>
#include <unordered_map>
#include <sstream>
#include <ashttp3lib/h1/utils.hpp>

namespace ashttp3lib::h1 {

//! \brief Request Class. Represents an HTTP/1.1 request.
class Request {
 public:
  std::string method;                    //!< HTTP method (GET, POST, etc.).
  std::string path;                      //!< Path specified in the request.
  std::unordered_map<std::string, std::string> headers;  //!< Request headers.
  std::string body;                      //!< Request body.

  //! \brief Constructor for Request class.
  //! \param request_stream. [std::istream&] Input string containing the HTTP request.
  Request(const std::string& requestString) noexcept {
    std::istringstream requestStream(requestString);

    // Parse the first line to get method, path, and HTTP version
    getline(requestStream, this->method, ' ');
    getline(requestStream, this->path, ' ');

    // Parse headers
    std::string line;
    while (getline(requestStream, line) && line != "\r") {
        size_t colonPos = line.find(':');
        if (colonPos != std::string::npos) {
            std::string headerName = line.substr(0, colonPos);
            // skip the colon after header name
            std::string headerValue = line.substr(colonPos + 2);  
            this->headers[headerName] = ashttp3lib::h1::utils::removeTrailingCarriageReturns(headerValue);
        }
    }

    // Parse the request body
    getline(requestStream, this->body, '\0');
  };
};

}  // namespace ashttp3lib::h1
// ashttp3lib/h1/server.hpp