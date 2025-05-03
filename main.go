package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	maxRange   = 10_000_000
	bufferSize = 1000
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <2-letter-country-code> [--v4|--v6]")
		return
	}

	cc := strings.ToUpper(os.Args[1])
	mode := os.Args[2]

	if len(cc) != 2 {
		fmt.Println("Error: Country code must be 2 letters (e.g. UA, US).")
		return
	}
	if mode != "--v4" && mode != "--v6" {
		fmt.Println("Error: You must specify either --v4 or --v6")
		return
	}

	useIPv6 := mode == "--v6"

	dbFile := "./db/IP2LOCATION-LITE-DB1.CSV"
	if useIPv6 {
		dbFile = "./db/IP2LOCATION-LITE-DB1.IPV6.CSV"
	}

	logFile, err := os.Create("skipped_" + cc + ".log")
	checkErr(err)
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	csvfile, err := os.Open(dbFile)
	checkErr(err)
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	outputFileName := cc + ".txt"
	if useIPv6 {
		outputFileName = cc + "_ipv6.txt"
	}

	file, err := os.Create(outputFileName)
	checkErr(err)
	defer file.Close()

	ipChan := make(chan string, bufferSize)
	var wg sync.WaitGroup

	go func() {
		for ip := range ipChan {
			_, err := file.WriteString(ip + "\n")
			checkErr(err)
		}
	}()

	var (
		totalIPs     uint64
		skippedCount int
		mutex        sync.Mutex
	)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Printf("CSV read error: %v", err)
			continue
		}

		if len(record) < 3 || record[2] != cc {
			continue
		}

		if useIPv6 {
			start := new(big.Int)
			end := new(big.Int)
			_, ok1 := start.SetString(record[0], 10)
			_, ok2 := end.SetString(record[1], 10)
			if !ok1 || !ok2 || start.Cmp(end) >= 0 {
				logger.Printf("Invalid IPv6 range: %s - %s", record[0], record[1])
				skippedCount++
				continue
			}

			count := new(big.Int).Sub(end, start)
			if count.Cmp(big.NewInt(maxRange)) > 0 {
				msg := fmt.Sprintf("Skipped large IPv6 range: %s IPs (%s - %s)",
					count.String(), bigIntToIPv6(start), bigIntToIPv6(end))
				fmt.Println(msg)
				logger.Println(msg)
				skippedCount++
				continue
			}

			wg.Add(1)
			go func(start, end *big.Int) {
				defer wg.Done()
				local := new(big.Int).Set(start)
				one := big.NewInt(1)
				var localCount uint64
				for local.Cmp(end) < 0 {
					ipChan <- bigIntToIPv6(local)
					local.Add(local, one)
					localCount++
				}
				mutex.Lock()
				totalIPs += localCount
				mutex.Unlock()
			}(new(big.Int).Set(start), new(big.Int).Set(end))

		} else {
			startIP, err1 := strconv.ParseUint(record[0], 10, 64)
			endIP, err2 := strconv.ParseUint(record[1], 10, 64)
			if err1 != nil || err2 != nil || startIP >= endIP {
				logger.Printf("Invalid IPv4 range: %v - %v", record[0], record[1])
				skippedCount++
				continue
			}

			count := endIP - startIP
			if count > maxRange {
				msg := fmt.Sprintf("Skipped large IPv4 range: %d IPs (%s - %s)", count, inet_ntoa(startIP), inet_ntoa(endIP))
				fmt.Println(msg)
				logger.Println(msg)
				skippedCount++
				continue
			}

			wg.Add(1)
			go func(start, end uint64) {
				defer wg.Done()
				var localCount uint64
				for ip := start; ip < end; ip++ {
					ipChan <- inet_ntoa(ip)
					localCount++
				}
				mutex.Lock()
				totalIPs += localCount
				mutex.Unlock()
			}(startIP, endIP)
		}
	}

	wg.Wait()
	close(ipChan)

	logger.Printf("TOTAL SKIPPED RANGES: %d\n", skippedCount)

	fmt.Printf("✅ Done. %d %s addresses saved to %s\n", totalIPs, func() string {
		if useIPv6 {
			return "IPv6"
		}
		return "IPv4"
	}(), outputFileName)
	fmt.Printf("📝 %d skipped/invalid ranges logged in skipped_%s.log\n", skippedCount, cc)
}

func inet_ntoa(ip uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func bigIntToIPv6(ip *big.Int) string {
	b := ip.FillBytes(make([]byte, 16))
	return net.IP(b).String()
}
