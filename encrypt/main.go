package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

var header []string
var ch = make(chan []string)

func main() {
	var batchSize int
	batchSize, _ = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if batchSize <= 0 {
		batchSize = 1000000
	}

	var key = []byte("example key 1234")
	go encryptedBatchWriter(batchSize, key)

	inFile, err := os.Open("data/bigdata.csv")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	csvReader := csv.NewReader(inFile)
	csvReader.Comment = '*'

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			close(ch)
			break
		}

		if header == nil {
			header = record
			ch <- header
			continue
		}

		ch <- record
	}
}

func encryptedBatchWriter(batchSize int, key []byte) {
	var batch = 1
	var count int

	f := newBatchWriter(key, batch)
	csvWriter := csv.NewWriter(f)

	for {
		select {
		case r, ok := <-ch:
			count++
			if count > batchSize {
				batch++
				count = 0

				csvWriter.Flush()
				err := f.Close()
				if err != nil {
					panic(err)
				}

				f = newBatchWriter(key, batch)
				csvWriter = csv.NewWriter(f)
				csvWriter.Write(header)
			}

			err := csvWriter.Write(r)
			if err != nil {
				panic(err)
			}

			if !ok {
				ch = nil
			}
		}

		if ch == nil {
			break
		}
	}

	csvWriter.Flush()
	err := f.Close()
	if err != nil {
		panic(err)
	}
}

func newBatchWriter(key []byte, batch int) *cipher.StreamWriter {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	outFile, err := os.OpenFile(fmt.Sprintf("output/batch%d.csv.bin", batch), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}

	return &cipher.StreamWriter{S: stream, W: outFile}
}
