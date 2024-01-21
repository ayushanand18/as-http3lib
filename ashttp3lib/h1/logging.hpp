/*
  ashttp3lib/h1/logging.hpp - A C++ HTTP/1.1 Library Logger Class
  
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

#include <chrono>
#include <ctime>
#include <iostream>
#include <sstream>

namespace ashttp3lib::logging {

//! \brief Logger Class. Provides logging functionality for different log levels.
class Logger {
 public:
  //! \brief Constructor for Logger class.
  Logger() {}

  //! \brief Destructor for Logger class.
  ~Logger() {}

  //! \brief Log an informational message.
  //! \param message. [const std::string&] The message to be logged.
  void info(const std::string& message) { log("INFO", message); }

  //! \brief Log a warning message.
  //! \param message. [const std::string&] The message to be logged.
  void warning(const std::string& message) { log("WARNING", message); }

  //! \brief Log an error message.
  //! \param message. [const std::string&] The message to be logged.
  void error(const std::string& message) { log("ERROR", message); }

 private:
  //! \brief Get the current timestamp in the format [%Y-%m-%d %H:%M:%S.%09d].
  //! \return [std::string] The current timestamp.
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

  //! \brief Log a message with a specified log level.
  //! \param level. [const std::string&] The log level (INFO, WARNING, ERROR).
  //! \param message. [const std::string&] The message to be logged.
  void log(const std::string& level, const std::string& message) {

    std::ostringstream logEntry;
    logEntry << getCurrentTimestamp() << " [" << level << "]: " << message
             << std::endl;
    std::cout << logEntry.str();  // Print to console as well
  }
};

}  // namespace ashttp3lib::logging

// ashttp3lib/h1/logging.hpp