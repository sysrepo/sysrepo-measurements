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
=======
* Architecture:          x86_64
* CPU op-mode(s):        32-bit, 64-bit
* Byte Order:            Little Endian
* CPU(s):                8
* On-line CPU(s) list:   0-7
* Thread(s) per core:    2
* Core(s) per socket:    4
* Socket(s):             1
* NUMA node(s):          1
* Vendor ID:             GenuineIntel
* CPU family:            6
* Model:                 94
* Model name:            Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
* Stepping:              3
* CPU MHz:               2500.000
* CPU max MHz:           2601.0000
* CPU min MHz:           800.0000
* BogoMIPS:              5183.92
* Virtualization:        VT-x
* L1d cache:             32K
* L1i cache:             32K
* L2 cache:              256K
* L3 cache:              6144K
* NUMA node0 CPU(s):     0-7
*
* Memory: 15G

Sysrepo was running on a local Docker container. It could also be useful to test
it on a device and remotely.


On ietf-interfaces module a number of most simple entries was added only
containing 'name' of interface and 'enable' value.

Tests were run few times each for combination of every operation and number of
entries in a request.

### Measures average amount of time for given number of entries and operation  processed in a single request
operation | number of entries | ms
--------- | ----------------- | --
delete 	1 	22.651006
delete 	2 	38.064887
delete 	4 	29.2917696667
delete 	8 	15.553244
delete 	16 	20.807633
delete 	32 	29.351339
delete 	64 	38.4092673333
delete 	128 	49.784648
delete 	256 	108.831738333
delete 	512 	373.296678
delete 	1024 	1397.51568633
delete 	2048 	5476.05380733
get 	1 	1.768185
get 	2 	2.171163
get 	4 	3.31087833333
get 	8 	5.640741
get 	16 	10.9224743333
get 	32 	24.0813563333
get 	64 	55.6608143333
get 	128 	149.439873
get 	256 	490.498112333
get 	512 	1818.744884
get 	1024 	6598.24484633
get 	2048 	26037.3821507
set 	1 	29.828248
set 	2 	16.202555
set 	4 	16.617822
set 	8 	26.491454
set 	16 	37.980052
set 	32 	27.7456426667
set 	64 	28.38218
set 	128 	48.7050303333
set 	256 	92.3466303333
set 	512 	177.184213
set 	1024 	454.813572
set 	2048 	1420.526593

![TestImg](/res/img/sysrepo_remote_perf_1.png)
