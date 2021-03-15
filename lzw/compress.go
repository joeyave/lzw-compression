package lzw

func Compress(str string) []int {
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
