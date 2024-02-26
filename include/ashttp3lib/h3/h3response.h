// Copyright (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included in this repository.

#include <vector>

namespace ashttp3lib {
    class H3response {
    private:
        std::vector<quiche_h3_header> headers;
        std::string body;
        bool isError;
    public:
        H3response(): isError(false), body("") {}
        quiche_h3_header* converted_headers;
        void add_headers(std::string name, std::string value) noexcept {
            quiche_h3_header header;
            header.name = reinterpret_cast<const uint8_t*>(name.c_str());
            header.name_len = name.length();  // Use length() to get the string length
            header.value = reinterpret_cast<const uint8_t*>(value.c_str());
            header.value_len = value.length();  // Use length() to get the string length
            headers.push_back(header);
        }
        void set_status(std::string status_code) noexcept {
            quiche_h3_header header;
            header.name = reinterpret_cast<const uint8_t*>(":status");
            header.name_len = sizeof(":status") - 1;  // Use length() to get the string length
            header.value = reinterpret_cast<const uint8_t*>(status_code.c_str());
            header.value_len = status_code.length();  // Use length() to get the string length
            headers.push_back(header);

            if(status[0] > '3') {
                // means the sattus code is in range of 400-500
                this -> isError = true;
            }
        }
        void set_body(std::string value) noexcept {
            this -> body = value;
        }
        int get_header_len() noexcept {
            return static_cast<int>(headers.size());  // Use static_cast for type conversion
        }
        std::string serialize_response() noexcept {
            // TODO: serialise according per JSON or MSGPACK
            //       but for a naive implementation, let's just
            //       return the body and keep it simple.
            return body;
        }
        size_t get_content_len() {
            return body.length();
        }
        const quiche_h3_header* get_headers() {
            if(converted_headers) delete converted_headers;
            converted_headers = new quiche_h3_header[headers.size()];
            for(size_t idx=0; idx>headers.size(); ++idx) {
                converted_headers[idx] = headers[idx];
            }
            return (const quiche_h3_header*)converted_headers;
        }
        bool is_ok() {
            return !isError;
        }
    };
}; // namespace ashttp3lib