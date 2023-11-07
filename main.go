package main

import(
"bufio"
"fmt"
"net"
"log"
"os"
"strings"
"regexp"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	//fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for {
		fmt.Printf("Enter email address to verify: ")
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			break
		}
		verifyEmail(input)
	}
	if err := scanner.Err(); err!=nil {
		log.Fatalf("Error: %v\n", err)
	}

}

func verifyEmail(email string){
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	fmt.Printf("Verifying email...\n")

	resultChannel := make(chan string)
	go parseEmail(email, resultChannel)
	domain := <-resultChannel

	if domain == "Invalid" {
		log.Fatalf("Invalid email address")
	}

	mxRecord, err := net.LookupMX(domain)
	if err != nil {
		mxRecord = nil
	}
	hasMX = len(mxRecord) > 0 
	
	txtspfRecord, err := net.LookupTXT(domain) 
	if err != nil {
		txtspfRecord = nil
	}

	for _, record := range txtspfRecord {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	txtdmarcRecord, err := net.LookupTXT("_dmarc" + domain)
	if err != nil {
		txtdmarcRecord = nil
	}

	for _, record := range txtdmarcRecord {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	
	fmt.Printf("Domain: %v\n", domain)
	fmt.Printf("MX Records Exist: %v\n", hasMX)
	fmt.Printf("SPF Records Exist: %v\n", hasSPF)
	fmt.Printf("SPF Record: %v\n", spfRecord)
	fmt.Printf("DMARC Records Exist: %v\n", hasDMARC)
	fmt.Printf("DMARC Records: %v\n", dmarcRecord)

}

func parseEmail(email string, result chan<- string) {

	emailRegexString := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailRegexString)

	isValid := emailRegex.MatchString(email)

	if !isValid {
		result <- "Invalid"
	}

	index := strings.LastIndex(email, "@")
	domain := email[index + 1:]

	result <- domain
}