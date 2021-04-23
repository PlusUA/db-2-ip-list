package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	startIP, finishIP uint64
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) == 2 {

		cc := os.Args[1]

		csvfile, err := os.Open("./db/IP2LOCATION-LITE-DB1.CSV")
		checkErr(err)

		r := csv.NewReader(csvfile)

		file, err := os.OpenFile(cc+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer file.Close()

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			checkErr(err)

			if record[2] == cc {

				if s1, err := strconv.ParseUint(record[0], 10, 64); err == nil {
					startIP = s1
				}
				if s2, err := strconv.ParseUint(record[1], 10, 64); err == nil {
					finishIP = s2
				}

				for i := startIP; i < finishIP; i++ {
					file.WriteString(inet_ntoa(i) + "\n")

				}
			}

		}
	} else {
		fmt.Println("Something wrong,country code please use 2 letter format as sample 'UA' for Ukraine or 'US' for USA")
	}
}

func inet_ntoa(ip uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8),
		byte(ip))
}
