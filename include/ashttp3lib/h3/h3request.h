// Copyright (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included in this repository.

namespace ashttp3lib {
    class H3request {
    private:
        quiche_h3_header* headers;
        std::string body;
    public:
        void set_headers(quiche_h3_header headers[], )
    }
} // namespace ashttp3lib