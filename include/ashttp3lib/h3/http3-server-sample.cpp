// Copytight (C) 2024, Ayush Anand
// This library has adopted many parts of Cloudflare's Quiche library
// and ethrefore we are including the license of the Quiche library
// in this file as well. The Quiche library is licensed under the
// BSD-2-Clause license and the license is included in this repository.

#include "http3-server.h"

int main() {
    ashttp3lib::Http3Server server("127.0.0.1", "8080");
    server.run();

    return 0;
}