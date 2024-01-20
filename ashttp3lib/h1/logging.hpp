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

  std::string getCurrentTimestamp() const {
    std::time_t now = std::time(0);
    struct std::tm* timeinfo = std::localtime(&now);
    char timestamp[20];
    std::strftime(timestamp, sizeof(timestamp), "[%Y-%m-%d %H:%M:%S]",
                  timeinfo);
    return std::string(timestamp);
  }

  void log(const std::string& level, const std::string& message) {
    
      std::ostringstream logEntry;
      logEntry << getCurrentTimestamp() << " [" << level << "]: " << message
               << std::endl;
      std::cout << logEntry.str();  // Print to console as well
  }
};

}  // namespace ashttp3lib::logging