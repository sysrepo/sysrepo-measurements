# About
This project is about implementing and describing some tests and measurement concerning [Sysrepo project](http://www.sysrepo.org).

## Test: multiple list entries in single request
Adapted example program from Sysrepo examples is used. Printouts on module
changes are removed not to interfere with basic time consumption.

Sysrepo was running on a local Docker container. It could also be useful to test it on a device and remotely.

Processor information:
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

| entries/ms | 	1	| 2	| 4	| 8	| 16 | 32 | 64 | 128 | 256 | 512 | 1024 | 2048 |
| ---------- | -----|---|---|---|----|----|----|-----|-----|-----|------|------|
| set | 29.828248 | 16.202555 | 16.617822 | 26.491454 | 37.980052 | 27.7456426667 | 28.38218 | 48.7050303333 | 92.3466303333 | 177.184213 | 454.813572 | 1420.526593 |
| get | 1.768185 | 2.171163 | 3.31087833333 | 5.640741 | 10.9224743333 | 24.0813563333 | 55.6608143333 | 149.439873 | 490.498112333 | 1818.744884 | 6598.24484633 | 26037.3821507 |
| delete | 22.651006 | 38.064887 | 29.2917696667 | 15.553244 | 20.807633 | 29.351339 | 38.4092673333 | 49.784648 | 108.831738333 | 373.296678 | 1397.51568633 | 5476.05380733 |

The set action run's the command.
```
<edit-config>
	<target><running/></target>
	<config>
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
			<interface>
				<name>eth0</name>
				<enabled>true</enabled>
			</interface>
		</interfaces>
	</config>
</edit-config>
```

The get action run's the command.
```
<get-config>
	<source><running/></source>
	<filter>
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
			<interface>
				<name>eth0</name>
				<enabled></enabled>
			</interface>
		</interfaces>
	</filter>
</get-config>
```

The delete action run's the command.
```
<edit-config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">
	<target><running/></target>
	<config>
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
			<interface nc:operation='delete'>
				<name>eth0</name>
			</interface>
		</interfaces>
	</config>
</edit-config>
```

It is observable that performance of getting leaves is deteriorating fastest
while setting is most stable operation and delete is in
between.
![TestImg](/res/img/sysrepo_remote_perf_1.png)
