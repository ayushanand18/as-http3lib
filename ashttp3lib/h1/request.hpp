#include <iostream>
#include <unordered_map>
#include "utils.hpp"

namespace ashttp3lib::h1 {
class Request {
 public:
  std::string method;
  std::string path;
  std::unordered_map<std::string, std::string> headers;
  std::string body;

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
      std::cerr << "Exception: " << e.what() << std::endl;
    }
  };
};
}  // namespace ashttp3lib::h1