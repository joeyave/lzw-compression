package lzw

import "fmt"

func Decompress(compressed []int) string {
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
