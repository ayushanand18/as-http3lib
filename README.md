# as-http3lib
> An Asynchronous HTTP/1.1 and HTTP/3 based Server Library with support of
> concurrency using ThreadPools and native Linux Primitives.

## Design
+ [Miro Board](https://miro.com/app/board/uXjVN645xLw=/?share_link_id=143978576535)

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

**Results**

Parameter        | ashttp3lib      | FastAPI
-----------------|-----------------|-----------------
Startup Time     | 13,008 ns       |
