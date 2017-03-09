package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sartura/go-netconf/netconf"
)

var mu = &sync.Mutex{}
var counter int

const numberOfConnections int = 8

func setTests(s *netconf.Session, wg *sync.WaitGroup, limit int) {

	setXML := `
	<edit-config>
		<target><running/></target>
		<config>
			<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
				<interface>
					<name>ethernet_%d</name>
					<type xmlns:ianaift="urn:ietf:params:xml:ns:yang:iana-if-type">ianaift:ethernetCsmacd</type>
				</interface>
			</interfaces>
		</config>
	</edit-config>`

	for {
		if counter >= limit {
			break
		}

		mu.Lock()
		item := counter
		counter++
		mu.Unlock()

		xml := fmt.Sprintf(setXML, item)

		_, err := s.Exec(netconf.RawMethod(xml))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			panic(err)
		}
	}

	wg.Done()
}

func main() {
	hostAddr := "0.0.0.0"
	user, passwd := "root", "root"

	// number of connections
	var s [numberOfConnections]*netconf.Session

	var err error
	for n := 0; n < numberOfConnections; n++ {
		s[n], err = netconf.DialSSH(hostAddr, netconf.SSHConfigPassword(user, passwd))
		if err != nil {
			log.Fatal(err)
		}
		defer s[n].Close()

		// need to make a netconf operation or netopeer-server will freeze after
		// 6 or more connections
		xml := "<commit/>"
		_, err := s[n].Exec(netconf.RawMethod(xml))
		if err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup

	fmt.Printf("\n\n\t\t%s", "set check")
	fmt.Printf("\n%-32s| %-15s | %-10s\n",
		"Operation", "number of items", "average time")

	start := time.Now()
	counter = 0
	//run tests with n connections
	for n := 0; n < numberOfConnections; n++ {
		wg.Add(1)
		go setTests(s[n], &wg, 10)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("%-32s| %-15d | %-10s\n", "set", 10, elapsed)

	start = time.Now()
	//run tests with n connections
	for n := 0; n < numberOfConnections; n++ {
		wg.Add(1)
		go setTests(s[n], &wg, 50)
	}
	wg.Wait()
	elapsed = time.Since(start)
	fmt.Printf("%-32s| %-15d | %-10s\n", "set", 50, elapsed)

	start = time.Now()
	//run tests with n connections
	for n := 0; n < numberOfConnections; n++ {
		wg.Add(1)
		go setTests(s[n], &wg, 100)
	}
	wg.Wait()
	elapsed = time.Since(start)
	fmt.Printf("%-32s| %-15d | %-10s\n", "set", 100, elapsed)
}
