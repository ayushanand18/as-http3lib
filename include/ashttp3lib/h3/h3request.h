// Copyright (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included in this repository.

#include <vector>
#include <string>

namespace ashttp3lib {
    class H3request {
    private:
        std::vector<quiche_h3_header> headers;
        std::string body;
        std::string path;
        std::string method;
    public:
        void add_headers(std::string name, std::string value) noexcept {
            quiche_h3_header header;
            header.name = reinterpret_cast<const uint8_t*>(name.c_str());
            header.name_len = name.length();  // Use length() to get the string length
            header.value = reinterpret_cast<const uint8_t*>(value.c_str());
            header.value_len = value.length();  // Use length() to get the string length
            if(name == ":path") {
                this->path = value;
            } else if (name == ":method") {
                this->method = method;
            } else {
                headers.push_back(header);
            }
        }
        void set_body(std::string value) noexcept {
            this->body = value;
        }
        int get_header_len() noexcept {
            return static_cast<int>(headers.size());  // Use static_cast for type conversion
        }
        std::string get_path() noexcept {
            return this->path;
        }
        std::string get_method() noexcept {
            return this->method;
        }
    };
} // namespace ashttp3lib
