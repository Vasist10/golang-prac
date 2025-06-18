package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcrecord \n")

	for scanner.Scan(){
		checkDomain(scanner.Text())

	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error could not read from input:\n", err )
	}

}

func checkDomain(domain string){
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	
	netRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
	}
	if len(netRecords) > 0 {
		hasMX = true
	} 

	txtRecords,err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v\n", domain, err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}
	dmarcTXTRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC TXT records for %s: %v\n", domain, err)
	}

	for _, record := range dmarcTXTRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	fmt.Printf("%s,%t,%t,%s,%t,%s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	
}

