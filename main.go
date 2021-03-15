package main

import (
	"bufio"
	"fmt"
	"github.com/joeyave/lzw-compression/lzw"
	"os"
)

func main() {
	fmt.Print("Input string to compress:\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()

	compressed := lzw.Compress(str)
	fmt.Println("\nCompressed:", compressed)

	decompressed := lzw.Decompress(compressed)
	fmt.Println("\nDecompressed:", decompressed)
}
