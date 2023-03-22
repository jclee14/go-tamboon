package main

import (
	"encoding/csv"
	"fmt"
	"go-tamboon/cipher"
	"io"
	"os"
	"runtime"
)

type transactionPayload struct {
	headers []string
	data    []string
	row     int
}

var doneCh = make(chan struct{})

func main() {

	// go func() {
	// 	for {
	// 		PrintMemUsage()
	// 		time.Sleep(100 * time.Nanosecond)
	// 	}
	// }()

	start()
}

func start() {
	// make chan
	job := make(chan transactionPayload, 10)

	workerAmount := 10
	for i := 0; i < workerAmount; i++ {
		go consumeAndCharge(job)
	}

	tmpFile, err := os.CreateTemp("", "tempFile.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	filePath := "data/fng.1000.csv.rot128"
	decryptFileAndWriteToFile(filePath, tmpFile)
	readAndProduceTransactionData(tmpFile, job)
}

func consumeAndCharge(ch <-chan transactionPayload) {
selectLoop:
	for {
		select {
		case payload := <-ch:
			fmt.Printf("%v\n", payload)
		case <-doneCh:
			break selectLoop
		}
	}
}

func decryptFileAndWriteToFile(filePath string, writeFile *os.File) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := cipher.NewRot128Reader(f)
	if err != nil {
		return err
	}

	chunkSize := 32
	for {
		chunk := make([]byte, chunkSize)
		n, err := r.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}
		// fmt.Printf("%s", string(chunk))

		if _, err := writeFile.Write(chunk[:n]); err != nil {
			return err
		}

		if err == io.EOF || n < chunkSize {
			break
		}
	}

	return nil
}

func readAndProduceTransactionData(tmpFile *os.File, ch chan<- transactionPayload) error {
	tmpFile.Seek(0, 0)
	csvReader := csv.NewReader(tmpFile)
	rowAmount := 0
	headers := make([]string, 0)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		rowAmount += 1

		if rowAmount == 1 {
			headers = append(headers, row...)
			continue
		}

		// do something with read line
		// fmt.Printf("%+v\n", row)

		// produce data to channel
		ch <- transactionPayload{
			headers: headers,
			data:    row,
			row:     rowAmount,
		}
	}

	doneCh <- struct{}{}

	fmt.Printf("Total row: %d\n", rowAmount)

	return nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v B", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v B", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v B", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b
}
