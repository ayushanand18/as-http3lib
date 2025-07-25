![as-http3lib](/assets/logo-startup.png)

# as-http3lib

[![Go Tests](https://github.com/ayushanand18/as-http3lib/actions/workflows/test-examples.yml/badge.svg)](https://github.com/ayushanand18/as-http3lib/actions/workflows/test-examples.yml)

An HTTP/1.1 + HTTP/2 + HTTP/3 server library written in pure Go. 

> [!NOTE]
> Faster alternative to FastAPI, natively written in Golang.

## Support
1. Supports HTTP v1.1, v2, v3.
2. Supports Self-signed TLS Certificates (guide to trust it locally on machine too).
3. Supports streaming HTTP responses (LLM/media use-cases).

## Table of Contents
+ [How to run](/notes/HOW-TO-RUN.md)
+ [See dev pipeline and next feature releases](/notes/TODO.md)
+ [Installation](#installation)
+ [Examples](/examples/)
+ [Performance and Benchmark results](/Performance.md)
+ [Acknowledgement](#acknowledgements)
+ [Documentation - Coming Soon!](/docs/)

## Installation
> Import the Go package in your service using
```sh
go mod add github.com/ayushanand18/as-http3lib
```

### Benchmark Results
> [!Note]
> More information [here](/notes/PERFORMANCE.md)

Parameter        | ashttp3lib::h1  | FastAPI (H/1.1)| ashttp3lib::h3  | ashttp3lib-go::h3 [latest]
-----------------|-----------------|----------------|-----------------|---------------------------
Startup Time     | 0.005 s         | 0.681 s        | 0.014 s         | 4.4499ms
RTT (p50)        |                 |                |                 | 1.751ms
RTT (p90)        | 6.88 ms         | 7.68 ms        | 4.49 ms         | 3.765ms
RTT (p95)        | 8.97 ms         | 9.34 ms        | 7.74 ms         | 4.796ms
RTT (p99)        |                 |                |                 | 7.678ms

> Tested by using `time` on Linux. These times are an average of 3 consecutive runs so as to
> offset system load irregularities however these figure might (and probably shall) differ on
> each and every run.

- Startup Time: 153x faster than FastAPI.
- 50.97% faster than FastAPI (p90).
- 48.65% faster than FastAPI (p95).

## Examples
This repository includes examples for building a simple server over [HTTP/3](./examples/naive/main.go).

## Acknowledgements
Other performant libraries form the backbone of this repository, and made it possible to build. We 
utilise the following open source libraries for developing `ashttp3lib`
- [cloudflare/quiche](https://github.com/cloudflare/quiche/) - older versions (in C++), not currently used
- [google/boringssl](https://github.com/google/boringssl/) - older versions (in C++), not currently used
- [quic-go](https://github.com/quic-go).

Thanks!
