#include <chrono>
#include <ctime>
#include <iostream>
#include <sstream>

namespace ashttp3lib::logging {

class Logger {
 public:
  Logger() {}

  ~Logger() {}

  void info(const std::string& message) { log("INFO", message); }

  void warning(const std::string& message) { log("WARNING", message); }

  void error(const std::string& message) { log("ERROR", message); }

 private:
  std::string getCurrentTimestamp() {
    auto now = std::chrono::system_clock::now();
    auto ns = std::chrono::duration_cast<std::chrono::nanoseconds>(
                  now.time_since_epoch())
                  .count();

    auto seconds =
        std::chrono::duration_cast<std::chrono::seconds>(now.time_since_epoch())
            .count();
    auto nanoseconds = ns - seconds * 1000000000;

    std::time_t time_seconds = static_cast<std::time_t>(seconds);
    struct std::tm* timeinfo = std::localtime(&time_seconds);

    char timestamp[30];  // Increased buffer size to accommodate nanoseconds
    std::strftime(timestamp, sizeof(timestamp), "[%Y-%m-%d %H:%M:%S]",
                  timeinfo);

    // Append nanoseconds to the timestamp
    char nsBuffer[10];  // 9 digits for nanoseconds
    std::snprintf(nsBuffer, sizeof(nsBuffer), ".%09d",
                  static_cast<int>(nanoseconds));

    return std::string(timestamp) + nsBuffer;
  }

  void log(const std::string& level, const std::string& message) {

    std::ostringstream logEntry;
    logEntry << getCurrentTimestamp() << " [" << level << "]: " << message
             << std::endl;
    std::cout << logEntry.str();  // Print to console as well
  }
};

}  // namespace ashttp3lib::logging