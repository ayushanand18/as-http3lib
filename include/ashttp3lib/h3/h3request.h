// Copyright (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included in this repository.

namespace ashttp3lib {
    class H3request {
    private:
        vector<quiche_h3_header> headers;
        std::string body;
    public:
        void add_headers(std::string name, std::string value) noexcept {
            headers.push_back({
                .name = (const uint8_t*)name;
                .name_len = sizeof(name) - 1;

                .value = (const uint8_t*)value;
                .value_len = sizeof(value) - 1;
            });
        }
        void set_body(std::string value) noexcept constexpr {
            this -> body = value;
        }
        int get_header_len() noexcept {
            return (int)headers.size();
        }
    }
} // namespace ashttp3lib