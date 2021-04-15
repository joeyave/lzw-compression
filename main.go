package main

import (
	"github.com/joeyave/lzw-compression/lzw"
	"log"
	"os"
	"path/filepath"
)

func main() {

	origFilePath, _ := filepath.Abs("test_data/orig/data.txt")
	fileToCompress, err := os.Open(origFilePath)
	if err != nil {
		panic(err)
	}

	compressedFilePath, _ := filepath.Abs("test_data/compr/data.txt")
	err = lzw.Compress(fileToCompress, compressedFilePath, 256)
	if err != nil {
		log.Fatal(err)
	}
	fileToCompress.Close()

	compressedFile, err := os.Open(compressedFilePath)
	if err != nil {
		panic(err)
	}

	decompressedFilePath, _ := filepath.Abs("test_data/decompr/data.txt")
	err = lzw.Decompress(compressedFile, decompressedFilePath, 256)
	if err != nil {
		log.Fatal(err)
	}
	compressedFile.Close()
}
