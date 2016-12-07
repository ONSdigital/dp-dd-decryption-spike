package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// 8 on a 4-core MBP seems optimum, a bit faster than 4, similar to 12
var concurrency = 8
var sem = make(chan int, concurrency)

func main() {
	var key = []byte("example key 1234")
	var wg sync.WaitGroup

	err := filepath.Walk("output", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		wg.Add(1)
		sem <- 1
		go func() {
			//log.Println(path)
			decryptFile(key, path)
			wg.Done()
			<-sem
		}()
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func decryptFile(key []byte, file string) {
	inFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Read the IV from the file first
	var iv [aes.BlockSize]byte
	_, err = inFile.Read(iv[:])
	if err != nil {
		panic(err)
	}
	stream := cipher.NewOFB(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: inFile}
	csvReader := csv.NewReader(reader)
	var header bool

	var c int

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		c++

		if !header {
			header = true
			continue
		}

		fmt.Fprintf(ioutil.Discard, "%s", record)

		// FIXME: maybe simulate DB write?
	}
	fmt.Fprintf(os.Stdout, "%d records\n", c)
}
