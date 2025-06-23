# as-http3lib
> An Asynchronous HTTP/1.1 and HTTP/3 based Server Library with support of
> concurrency using Event Loop (using `libev`) and native Linux Primitives.

## Installation
> Currently the compiled static libraries support only Linux based distros (tested only on Ubuntu 22.04).
* Before installing the project, make sure you have [`libev`](https://github.com/enki/libev) installed on your machine.
* Clone this repository on your machine. ANd you are good to go!

This project uses [quiche](https://github.com/cloudflare/quiche) which also uses [boringssl](github.com/google/boringssl). 
Those who want to go with their own build of static libraries, can build these two libraries from source and add it to `deps/`.

## Performance
Some stats:
* Minimal version 138x (138 times) faster on startup time than FastAPI based server (equal endpoints and configs).

## Benchmarks
During the development of this project, we ran a few benchmarks to evaluate the 
performance of this library against a few popular libraries. The results are illustrated
in the table below.

**Machine Information**
```
Boost: v1.84.0

g++: g++ (Ubuntu 9.4.0-1ubuntu1~20.04.2) 9.4.0
Copyright (C) 2019 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

compilation flags: --std=c++2a -Wall -Wextra -pedantic -lboost_system -pthread

make: GNU Make 4.2.1
Built for x86_64-pc-linux-gnu
Copyright (C) 1988-2016 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Architecture:                       x86_64
CPU op-mode(s):                     32-bit, 64-bit
Byte Order:                         Little Endian
Address sizes:                      48 bits physical, 48 bits virtual
CPU(s):                             2
On-line CPU(s) list:                0,1
Thread(s) per core:                 2
Core(s) per socket:                 1
Socket(s):                          1
NUMA node(s):                       1
Vendor ID:                          AuthenticAMD
CPU family:                         25
Model:                              1
Model name:                         AMD EPYC 7763 64-Core Processor
Stepping:                           1
CPU MHz:                            3239.193
BogoMIPS:                           4890.86
Virtualization:                     AMD-V
Hypervisor vendor:                  Microsoft
Virtualization type:                full
L1d cache:                          32 KiB
L1i cache:                          32 KiB
L2 cache:                           512 KiB
L3 cache:                           32 MiB
NUMA node0 CPU(s):                  0,1
```

### Results

Parameter        | ashttp3lib::h1  | FastAPI (H/1.1)| ashttp3lib::h3  | ashttp3lib-go::h3
-----------------|-----------------|----------------|-----------------|-------------------
Startup Time     | 0.005 s         | 0.681 s        | 0.014 s         | 20ms
RRT (p50)        |                 |                |                 | 1.995ms
RRT (p90)        | 6.88 ms         | 7.68 ms        | 4.49 ms         | 4.497ms
RRT (p95)        | 8.97 ms         | 9.34 ms        | 7.74 ms         | 5.837ms
RRT (p99)        |                 |                |                 | 11.13ms

> Tested by using `time` on Linux. These times are an average of 3 consecutive runs so as to
> offset system load irregularities however these figure might (and probably shall) differ on
> each and every run.

- 41.54% faster than FastAPI (p90).
- 17.13% faster than FastAPI (p95).

## Examples
This repository includes examples for building a simple server over [HTTP/1.1](./examples/h1_server.cpp) and [HTTP/3](./examples/http3-server-sample.cpp).

## Acknowledgements
Other performant libraries form the backbone of this repository, and made it possible to build. We 
utilise the following open source libraries for developing `ashttp3lib`
- [cloudflare/quiche](https://github.com/cloudflare/quiche/)
- [google/boringssl](https://github.com/google/boringssl/)
- [libev](http://software.schmorp.de/pkg/libev.html)
- [uthash](https://troydhanson.github.io/uthash/)

Thanks!
