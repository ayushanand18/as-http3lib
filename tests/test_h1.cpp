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

// Test case for request parsing by ashttp3lib::h1::request::Request Class, GET
TEST(AsHttp3LibTest, GETRequestParsing) {
    std::string input_stream = "GET /test HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Android Gecko\r\nAccept: */*\r\n\r\n";
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

// Test case for request parsing by ashttp3lib::h1::request::Request Class, POST
TEST(AsHttp3LibTest, POSTRequestParsing) {
    std::string input_stream = "POST /test HTTP/1.1\r\n"
        "Host: example.com\r\n"
        "User-Agent: Mozilla/5.0\r\n"
        "Content-Type: application/json\r\n"
        "Content-Length: 19\r\n"
        "\r\n"
        "{\"key\":\"value\"}";
    ashttp3lib::h1::Request result(input_stream);

    // Assert, all headers are parsed directly
    ASSERT_EQ(result.headers.size(), 4);
    // Assert, method is as passed
    ASSERT_EQ(result.method, "POST");
    // Assert, path is as passed
    ASSERT_EQ(result.path, "/test");
    // Assert, pick a random header 
    ASSERT_EQ(result.headers.at("Host"), "example.com");
    // Assert, the body is correctly parsed
    ASSERT_EQ(result.body, "{\"key\":\"value\"}");
}
