#include "../ashttp3lib/h1/server.hpp"

std::string handleGetRequest(ashttp3lib::h1::Request& request_packet) {
    return "Handling POST request";
}

int main() {
  try {
    ashttp3lib::h1::HTTPServer server(8002);
    server.get("/test", handleGetRequest);
    server.run();
  } catch (const std::exception& e) {
    std::cerr << "Exception: " << e.what() << std::endl;
  }

  return 0;
}