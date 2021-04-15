package main

import (
	"github.com/joeyave/lzw-compression/lzw"
	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("test_data/orig/data.txt")
	if err != nil {
		log.Fatal(err)
	}

	compressed := lzw.Compress(data)
	ioutil.WriteFile("test_data/compr/data.txt", compressed, 0666)

	decompressed := lzw.Decompress(compressed)
	ioutil.WriteFile("test_data/decompr/data.txt", decompressed, 0666)
}
