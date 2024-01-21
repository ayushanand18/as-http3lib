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
#include "utils.hpp"

namespace ashttp3lib::h1 {

//! \brief Request Class. Represents an HTTP/1.1 request.
class Request {
 public:
  std::string method;                    //!< HTTP method (GET, POST, etc.).
  std::string path;                      //!< Path specified in the request.
  std::unordered_map<std::string, std::string> headers;  //!< Request headers.
  std::string body;                      //!< Request body.

  //! \brief Constructor for Request class.
  //! \param request_stream. [std::istream&] Input stream containing the HTTP request.
  Request(std::istream& request_stream) noexcept {
    request_stream.exceptions(std::istream::failbit|std::istream::badbit);
    std::string line = "";
    try {
      for (int line_no = 0; getline(request_stream, line); ++line_no) {
        if(request_stream.fail()) {
          request_stream.clear();
          break;
        }
        if (line_no == 0) {
          auto splitted_title = ashttp3lib::h1::utils::split(line, " ");
          if(splitted_title.size() < 2) break;
          this->method = splitted_title[0];
          this->path = splitted_title[1];
        } else {
          auto splitted_line = ashttp3lib::h1::utils::split(line, ": ");
          if(splitted_line.size() < 2) break;
          this->headers[splitted_line[0]] = splitted_line[1];
        }
      }
    } catch (const std::exception& e) {
      std::cerr << "Exception: " << e.what() << "\n";
    }
  };
};

}  // namespace ashttp3lib::h1
// ashttp3lib/h1/server.hpp