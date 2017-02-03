# About
This project is about implementing and describing some tests and measurement concerning [Sysrepo project](http://www.sysrepo.org).

## Test: multiple list entries in single request
Adapted example program from Sysrepo examples is used. Printouts on module
changes are removed not to interfere with basic time consumption.

Processor:
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                8
On-line CPU(s) list:   0-7
Thread(s) per core:    2
Core(s) per socket:    4
Socket(s):             1
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 94
Model name:            Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
Stepping:              3
CPU MHz:               2500.000
CPU max MHz:           2601.0000
CPU min MHz:           800.0000
BogoMIPS:              5183.92
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              6144K
NUMA node0 CPU(s):     0-7

Memory: 15G

Sysrepo was running on a local Docker container. It could also be useful to test it on a device and remotely. 

![TestImg](/res/img/sysrepo-remote-perf_1.png)
