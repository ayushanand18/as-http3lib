#include <gtest/gtest.h>
#include "../ashttp3lib/h1/request.hpp"

// Test case for ashttp3lib::h1::utils::split() function
TEST(AsHttp3LibTest, SplitFunction) {
    std::vector<std::string> run_result;
    std::vector<std::string> expected_result {"hello", "world"};

    run_result = ashttp3lib::h1::utils::split("hello world", " ");

    // Assert
    ASSERT_EQ(run_result, expected_result);
}

// Test case for request parsing by ashttp3lib::h1::request::Request Class
TEST(AsHttp3LibTest, RequestParsing) {
    std::string input_stream = "GET /test HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Your-User-Agent\r\nAccept: */*\r\n\r\n";
    ashttp3lib::h1::Request result(input_stream);

    // Assert, all headers are parsed directly
    ASSERT_EQ(result.headers.size(), 3);
    // Assert, method is as passed
    ASSERT_EQ(result.method, "GET");
    // Assert, path is as passed
    ASSERT_EQ(result.path, "/test");
    // Assert, pick a random header 
    ASSERT_EQ(result.headers.at("Host"), "example.com");
}
