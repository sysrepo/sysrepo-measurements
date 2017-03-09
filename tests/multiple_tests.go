package main

import (
	"bytes"
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
	defer s.Close()

	_, err = s.Exec(netconf.RawMethod(xml))
	if err != nil {
		fmt.Printf("delete data ERROR: %s\n", err)
	}
}

func fillDatastore(numberOfItems int) {
	counter = numberOfItems

	if numberOfItems == 0 {
		return
	}

	leftXML := `
<edit-config>
	<target><running/></target>
	<config xmlns:op="urn:ietf:params:xml:ns:netconf:base:1.0">
		<interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">`

	middleXML := `
			<interface op:operation="create">
				<name>ethernet_%d</name>
				<type xmlns:ianaift="urn:ietf:params:xml:ns:yang:iana-if-type">ianaift:ethernetCsmacd</type>
			</interface>`

	rightXML := `
		</interfaces>
	</config>
</edit-config>`

	var buffer bytes.Buffer

	buffer.Write([]byte(leftXML))
	for i := 0; i < numberOfItems; i++ {
		buffer.Write([]byte(fmt.Sprintf(middleXML, i)))
	}
	buffer.Write([]byte(rightXML))

	s, err := netconf.DialSSH(hostAddr, netconf.SSHConfigPassword(user, passwd))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	_, err = s.Exec(netconf.RawMethod(buffer.String()))
	if err != nil {
		fmt.Printf("init data ERROR: %s\n", err)
	}
}

func setTests(s *netconf.Session, wg *sync.WaitGroup, limit int) {
	setXML := `
<edit-config>
	<target><running/></target>
	<config xmlns:op="urn:ietf:params:xml:ns:netconf:base:1.0">
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
			fmt.Printf("set item ERROR: %s\n", err)
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
	sessions := []int{1, 4}
	elements := []int{10, 100}
	existingItems := []int{0, 1000}

	var wg sync.WaitGroup
	counter = 0

	for _, existingItem := range existingItems {
		for _, numberOfSessions := range sessions {
			fmt.Printf("\n\n\tset itm with  %d connections and %d existing items", numberOfSessions, existingItem)
			fmt.Printf("\n%-32s| %-15s | %-10s\n", "Operation", "number of items", "total time")
			fmt.Printf("-------------------------------------------------------------------\n")
			for _, element := range elements {
				cleanDatastore()
				fillDatastore(existingItem)
				sess := getSessions(numberOfSessions)
				defer closeSessions(sess)

				start := time.Now()
				for _, s := range sess {
					wg.Add(1)
					go setTests(s, &wg, element+counter)
				}
				wg.Wait()
				elapsed := time.Since(start)
				fmt.Printf("%-32s| %-15d | %-10s\n", "set", element, elapsed)
			}
		}
	}
}
