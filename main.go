package main

import (
	"github.com/joeyave/lzw-compression/lzw"
	"log"
	"os"
)

func main() {
	input, err := os.Open("test_data/orig/data.txt")
	if err != nil {
		panic(err)
	}
	output, err := os.Create("test_data/compressed/data.txt")
	if err != nil {
		panic(err)
	}

	table := make(map[string]int64)
	for code := 0; code < 256; code++ {
		table[string(rune(code))] = int64(code)
	}

	conf, _ := lzw.SetupConfig(input, output, 20, table)
	if opt == "-c" {
		err := lzw.Compress(conf)
		if err != nil {
			log.Fatal(err)
		}
	} else if opt == "-d" {
		err := lzw.Decompress(conf)
		if err != nil {
			log.Fatal(err)
		}
	}
}
