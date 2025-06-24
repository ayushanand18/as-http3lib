# Design
> Design Document containing design of the library and class.

## Classes
```c++
class ClientConnection {
  UUID request-id;
  IP destination_ip;
  Connection connection;
  queue<Request> request_queue;
};

class ResponsePacket {
  Response http_response;
  IP destination_ip;
  Time timestamp;
};

class Response {
  Header headers;
  StatusCode status;
  string body;
  MSG_PACK(); // serialiser
};

class Request {
  Protocol protocol;
  Header headers;
  string body;
  MSG_UNPACK() // de-serialiser
};

class HttpServer {
  int portnumber;
  IP ip_address;
  vector<Connection> active_connections;
  ThreadPoolExecutor threadpool;
  unordered_map<pair<string, string>, function* () -> T> etf_binder; // endpoint to function binder
      [[protocol, endpoint] -> function *]
  queue<Response> response_queue;
  EventLoop eventloop;

  // methods
  get(string endpoint, function* bind_function) -> any // to bind functions to listen to GET method
  post(string endpoint, function* bind_function) -> any // to bind functions to listen to POST method
  [...] more HTTP methods
};
```
