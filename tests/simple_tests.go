package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sartura/go-netconf/netconf"
)

const (
	LOOP_COUNT = 1
	req1       = `
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
	</get-config>`

	set1 = `
	<edit-config>
		<target><running/></target>
		<config>
			<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
				<interface>
					<name>%s</name>
					<enabled>%v</enabled>
				</interface>
			</interfaces>
		</config>
	</edit-config>`

	set2 = `
	<edit-config>
		<target><running/></target>
		<config>
			<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
				<interface>
					<name>%s</name>
					<enabled>%v</enabled>
				</interface>
				<interface>
					<name>%s</name>
					<enabled>%v</enabled>
				</interface>
			</interfaces>
		</config>
	</edit-config>`

	del1 = `
	<edit-config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0">
		<target><running/></target>
		<config>
			<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
				<interface nc:operation='delete'>
					<name>%s</name>
				</interface>
			</interfaces>
		</config>
	</edit-config>`
)

func generateSetRequest(nEntires int) string {

	magic := 42
	n := nEntires + magic
	preReq := `<edit-config><target><running/></target><config><interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">`
	postReq := `</interfaces></config></edit-config>`
	entryFmt := `<interface><name>%s</name><type xmlns:ianaift="urn:ietf:params:xml:ns:yang:iana-if-type">ianaift:ethernetCsmacd</type></interface>`

	req := "" + preReq

	for i := magic; i < n; i++ {
		entry := fmt.Sprintf(entryFmt, "eth"+strconv.Itoa(i))
		req = req + entry
	}

	return req + postReq
}

func generateGetRequest(nEntires int) string {

	magic := 42
	n := nEntires + magic

	preReq := `<get-config><source><running/></source><filter><interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">`
	postReq := `</interfaces></filter></get-config>`
	entryFmt := `<interface><name>%s</name><enabled></enabled></interface>`

	req := "" + preReq

	for i := magic; i < n; i++ {
		entry := fmt.Sprintf(entryFmt, "eth"+strconv.Itoa(i))
		req = req + entry
	}

	return req + postReq
}

func generateDeleteRequest(nEntires int) string {

	magic := 42
	n := nEntires + magic

	preReq := `<edit-config xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0"><target><running/></target><config><interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">`
	postReq := `</interfaces></config></edit-config>`
	entryFmt := `<interface nc:operation='delete'><name>%s</name></interface>`

	req := "" + preReq

	for i := magic; i < n; i++ {
		entry := fmt.Sprintf(entryFmt, "eth"+strconv.Itoa(i))
		req = req + entry
	}

	return req + postReq
}

// This test sets node to some value and commits it LOOP_COUNT number of times.
func testSet(s *netconf.Session) {

	var setReq string
	setReqTrue := fmt.Sprintf(set1, "eth0", true)
	setReqFalse := fmt.Sprintf(set1, "eth0", false)

	start := time.Now()
	for i := 0; i < LOOP_COUNT; i++ {
		fmt.Print(".")

		if i%2 == 0 {
			setReq = setReqFalse
		} else {
			setReq = setReqTrue
		}

		_, err := s.Exec(netconf.RawMethod(setReq))
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("[testSet] Elapsed time:", elapsed)

	fmt.Println(setReq)
}

// This test sets and gets node to some value and commits it LOOP_COUNT number of times.
func testSetGet(s *netconf.Session) {

	var setReq string
	setReqTrue := fmt.Sprintf(set1, "eth0", true)
	setReqFalse := fmt.Sprintf(set1, "eth0", false)
	getReq := req1

	start := time.Now()
	for i := 0; i < LOOP_COUNT; i++ {
		fmt.Print(".")

		if i%2 == 0 {
			setReq = setReqFalse
		} else {
			setReq = setReqTrue
		}

		// fmt.Println(getReq)
		_, err := s.Exec(netconf.RawMethod(setReq))
		if err != nil {
			panic(err)
		}

		_, err = s.Exec(netconf.RawMethod(getReq))
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("[testSetGet] Elapsed time:", elapsed)

	fmt.Println(setReq)
}

func testSetDelete(s *netconf.Session) {
	var setReq string
	setReqTrue := fmt.Sprintf(set2, "eth0", true, "eth1", true)
	setReqFalse := fmt.Sprintf(set2, "eth0", false, "eth1", true)
	delReq := fmt.Sprintf(del1, "eth0")

	start := time.Now()
	for i := 0; i < LOOP_COUNT; i++ {
		fmt.Print(".")

		if i%2 == 0 {
			setReq = setReqFalse
		} else {
			setReq = setReqTrue
		}

		_, err := s.Exec(netconf.RawMethod(setReq))
		if err != nil {
			panic(err)
		}

		_, err = s.Exec(netconf.RawMethod(delReq))
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("[testSetDelete] Elapsed time:", elapsed)
}

func testSetNTimes(s *netconf.Session, n int) time.Duration {

	setReq := generateSetRequest(n)

	start := time.Now()

	_, err := s.Exec(netconf.RawMethod(setReq))
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Println("set", n, ":", elapsed)

	return elapsed
}

func testGetNTimes(s *netconf.Session, n int) time.Duration {

	getReq := generateGetRequest(n)

	start := time.Now()

	_, err := s.Exec(netconf.RawMethod(getReq))
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Println("get", n, ":", elapsed)

	return elapsed
}

func testDeleteNTimes(s *netconf.Session, n int) time.Duration {

	delReq := generateDeleteRequest(n)

	start := time.Now()

	_, err := s.Exec(netconf.RawMethod(delReq))
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Println("delete", n, ":", elapsed)

	return elapsed
}

// Test description => size of data => measured times
type ReportMap map[string]map[int][]time.Duration

// Compute average mean of time values.
func avgTime(ts []time.Duration) time.Duration {
	var totalTime time.Duration

	for _, t := range ts {
		totalTime = totalTime + t
	}

	return totalTime / time.Duration(len(ts))
}

func testScalability(s *netconf.Session) {

	iterCnt := 1
	logCnt := 1

	r := make(ReportMap)
	r["set"] = make(map[int][]time.Duration)
	r["get"] = make(map[int][]time.Duration)
	r["delete"] = make(map[int][]time.Duration)

	i := 1

	for j := 0; j < logCnt; j++ {
		r["set"][i] = make([]time.Duration, iterCnt)
		r["get"][i] = make([]time.Duration, iterCnt)
		r["delete"][i] = make([]time.Duration, iterCnt)

		i = i * 2
	}

	for n := 0; n < iterCnt; n++ {

		i := 1

		for j := 0; j < logCnt; j++ {
			r["set"][i][n] = r["set"][i][n] + testSetNTimes(s, i)
			r["get"][i][n] = r["get"][i][n] + testGetNTimes(s, i)
			r["delete"][i][n] = r["delete"][i][n] + testDeleteNTimes(s, i)

			i = i * 2
		}
	}

	printReport(r)
	dumpReport(r, "./report.txt")
}

func printMeasureHeader(title string) {

	fmt.Printf("\n\n\t\t%s", title)
	fmt.Printf("\n%-32s| %-15s | %-10s\n",
		"Operation", "number of items", "average time")
	fmt.Printf("---------------------------------------------------------------------------------------------------\n")
}

func printReport(r ReportMap) {

	printMeasureHeader(("Scalability check"))

	for description, sizeTimes := range r {
		for size_, times := range sizeTimes {
			size := size_
			fmt.Printf("%-32s| %-15d | %-10s\n", description, size, avgTime((times)).String())

		}
	}
}

func dumpReport(r ReportMap, filePath string) {

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for description, sizeTimes := range r {
		for size_, times := range sizeTimes {
			size := size_
			fmt.Fprintf(f, "%s,%d,%s\n", description, size, times)

		}
	}
}

func main() {
	hostAddr := "0.0.0.0:10830"
	user, passwd := "root", "root"

	sessionStart := time.Now()
	s, err := netconf.DialSSH(hostAddr, netconf.SSHConfigPassword(user, passwd))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	sessionEstablished := time.Since(sessionStart)
	fmt.Println("Establishment of session took", sessionEstablished)

	testScalability(s)
}
