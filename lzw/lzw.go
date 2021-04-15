package lzw

import (
	"fmt"
	"strconv"
)

const mapSize int64 = 256

func Compress(data []byte) []byte {
	if len(data) < 1 {
		return []byte{}
	}

	code := mapSize
	dictionary := make(map[string]int64, code)

	for i := int64(0); i < mapSize; i++ {
		dictionary[string(i)] = i
	}

	result := make([]int64, 0)
	word := make([]byte, 0)

	for _, c := range data {
		phrase := append(word, c)

		_, ok := dictionary[string(phrase)]
		if ok {
			word = phrase
		} else {
			result = append(result, dictionary[string(word)])
			dictionary[string(phrase)] = code
			code++
			word = []byte{c}
		}
	}

	if len(word) > 0 {
		result = append(result, dictionary[string(word)])
	}

	return toBytes(result)
}

func Decompress(data []byte) []byte {

	if len(data) < 1 {
		return []byte{}
	}

	compressed := fromBytes(data)

	code := mapSize
	dictionary := make(map[int64][]byte, code)

	for i := int64(0); i < mapSize; i++ {
		dictionary[i] = []byte{byte(i)}
	}

	result := make([]byte, 0)
	word := make([]byte, 0)

	for _, k := range compressed {
		var entry []byte
		x, ok := dictionary[k]

		if ok {
			entry = x[:len(x):len(x)]
		} else if k == code && len(word) > 0 {
			entry = append(word, word[0])
		} else {
			return []byte{}
		}
		result = append(result, entry...)

		if len(word) > 0 {
			word = append(word, entry[0])
			dictionary[code] = word
			code++
		}

		word = entry
	}

	return result
}

func toBytes(data []int64) []byte {
	bitsStr := ""
	bits := 9
	idx := mapSize
	for _, v := range data {
		if idx > 1<<bits {
			bits++
		}
		format := "%0" + strconv.Itoa(bits) + "b"
		bitsStr += fmt.Sprintf(format, v)
		idx++
	}

	bytes := make([]byte, 0)
	length := len(bitsStr)
	size := 8
	for i := 0; i < length; i += size {
		end := i + size
		if end > length {
			end = length
		}
		bitStr := bitsStr[i:end]
		bit, err := strconv.ParseInt(bitStr, 2, size+1)
		if err != nil {
			fmt.Println(err.Error())
		}
		bytes = append(bytes, byte(bit))
	}

	return bytes
}

func trimByteStr(str string) string {
	length := len(str)
	trim := ""
	idx := 8
	bits := 9
	i := 0
	for ; i < length; i += bits {
		if i+bits > length {
			break
		}
		if idx > 1<<bits {
			bits++
		}
		idx++
	}

	size := 8
	mod := (length - i) % size
	get := size - mod
	elem := str[length-get : length]
	trim = str[:length-size] + elem

	return trim
}

func fromBytes(data []byte) []int64 {
	byteStr := ""
	for _, b := range data {
		bitStr := fmt.Sprintf("%08b", b)
		byteStr += bitStr
	}

	byteStr = trimByteStr(byteStr)

	compressed := make([]int64, 0)
	bits := 9
	length := len(byteStr)
	idx := 256
	for i := 0; i < length; i += bits {
		if idx > 1<<bits {
			bits++
		}
		start := i
		end := i + bits
		if end > length {
			end = length
			start = end - bits
		}
		bitStr := byteStr[start:end]
		bit, err := strconv.ParseInt(bitStr, 2, bits+1)
		if err != nil {
			fmt.Println(err.Error())
		}
		compressed = append(compressed, bit)
		idx++
	}
	return compressed
}
