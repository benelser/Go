package main

import(
	"bufio"
	"fmt"
	"os"
	"log"
	"regexp"
)

func main(){
	ipAdresses := make(map[string]int)
	valueCounter := 1
	csvFile, _ := os.Open("C:\\SOMEFOLDER\\IpTest.csv")
	defer csvFile.Close()
	scanner := bufio.NewScanner(csvFile)
	scanner.Split(bufio.ScanLines)
	for{
		// Read to next token...In this case line
		line := scanner.Scan()
		if line == false {
			// returns false on error or EOF check err
			err := scanner.Err()
			if err == nil {
				log.Println("Scan reached EOF")
				break
			} else {
				log.Fatal(err)
			}
		}
		
		// Get data from scan with Bytes() or Text()
		regEx := regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
		matches := regEx.FindAll(scanner.Bytes(), -1)
		for _, match := range matches {
			if _, ok := ipAdresses[string(match)]; ok {
				continue
			}
			ipAdresses[string(match)] = valueCounter
			valueCounter ++
		}
	}

	for key, value := range ipAdresses {
		fmt.Printf("Key: %v Value:%v\r\n", key, value)
	}
}