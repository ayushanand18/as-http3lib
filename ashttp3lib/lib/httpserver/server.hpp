#include <cstdint>

//! \brief Implements a msgpack-h/3 server. This is the main interfacing
//! point with the library for creating servers.
//!
//! The server maintains a registry of function bindings that it uses to
//! dispatch calls. It also takes care of managing worker threads and QUIC
//! connections.
//! The server does not start listening right after construction in order
//! to allow binding functions before that. Use the `run` or `async_run`
//! functions to start listening on the port.
//! This class is not copyable, but moveable.
namespace ashttp3lib :: httpserver {
    class HttpServer {
        //! \brief Constructs a server that listens on the localhost on the
        //! specified port.
        //!
        //! \param port The port number to listen on.
        explicit HttpServer(uint16_t port);

        //! \brief Move constructor. This is implemented by calling the
        //! move assignment operator.
        //!
        //! \param other The other instance to move from.
        HttpServer(HttpServer&& other) noexcept;
    };
}
