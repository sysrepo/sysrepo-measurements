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

const hostAddr string = "0.0.0.0"
const user string = "root"
const passwd string = "root"

func cleanDatastore() {
	counter = 0

	xml := `
<edit-config>
	<target><running/></target>
	<config xmlns:op="urn:ietf:params:xml:ns:netconf:base:1.0">
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces" op:operation='delete'/>
	</config>
</edit-config>`

	s, err := netconf.DialSSH(hostAddr, netconf.SSHConfigPassword(user, passwd))
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Exec(netconf.RawMethod(xml))
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

func setTests(s *netconf.Session, wg *sync.WaitGroup, limit int) {

	setXML := `
<edit-config>
	<target><running/></target>
	<config xmlns:op="urn:ietf:params:xml:ns:netconf:base:1.0">
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
			<interface op:operation="create">
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
		}
	}

	wg.Done()
}

func getSessions(numberOfSessions int) []*netconf.Session {
	// number of connections
	var s []*netconf.Session

	for n := 0; n < numberOfSessions; n++ {
		session, err := netconf.DialSSH(hostAddr, netconf.SSHConfigPassword(user, passwd))
		if err != nil {
			log.Fatal(err)
		}
		s = append(s, session)
	}

	return s
}

func closeSessions(sessions []*netconf.Session) {
	for _, s := range sessions {
		s.Close()
	}
}

func main() {
	sessions := []int{1, 2, 4}
	limits := []int{10, 20, 40}

	var wg sync.WaitGroup
	counter = 0

	for _, numberOfSessions := range sessions {
		sess := getSessions(numberOfSessions)
		fmt.Printf("\n\n\t\tset check, with  %d connections", numberOfSessions)
		fmt.Printf("\n%-32s| %-15s | %-10s\n", "Operation", "number of items", "total time")
		fmt.Printf("-------------------------------------------------------------------\n")
		for _, limit := range limits {
			cleanDatastore()
			start := time.Now()
			for _, s := range sess {
				wg.Add(1)
				go setTests(s, &wg, limit)
			}
			wg.Wait()
			elapsed := time.Since(start)
			fmt.Printf("%-32s| %-15d | %-10s\n", "set", limit, elapsed)
		}
		closeSessions(sess)
	}
}
