package main

import (
	"bufio"
	"fmt"
	"os"
)

func compressLZW(str string) []int {
	code := 256
	dictionary := make(map[string]int)

	for i := 0; i < 256; i++ {
		dictionary[string(i)] = i
	}

	currChar := ""
	result := make([]int, 0)

	for _, c := range []byte(str) {
		phrase := currChar + string(c)

		_, ok := dictionary[phrase]
		if ok {
			currChar = phrase
		} else {
			result = append(result, dictionary[currChar])
			dictionary[phrase] = code
			code++
			currChar = string(c)
		}
	}

	if currChar != "" {
		result = append(result, dictionary[currChar])
	}

	return result
}

func decompressLZW(compressed []int) string {
	code := 256
	dictionary := make(map[int]string)
	for i := 0; i < 256; i++ {
		dictionary[i] = string(i)
	}

	currChar := string(compressed[0])
	result := currChar

	for _, element := range compressed[1:] {
		var word string

		x, ok := dictionary[element]
		if ok {
			word = x
		} else if element == code {
			word = currChar + currChar[:1]
		} else {
			panic(fmt.Sprintf("Bad compressed element: %d", element))
		}

		result += word

		dictionary[code] = currChar + word[:1]
		code++

		currChar = word
	}

	return result
}

func main() {
	fmt.Print("Input string to compress:\n")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	str := scanner.Text()

	compressed := compressLZW(str)
	fmt.Println("\nCompressed:", compressed)

	decompressed := decompressLZW(compressed)
	fmt.Println("\nDecompressed:", decompressed)
}
