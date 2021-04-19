package lzw

import (
	"github.com/dgryski/go-bitstream"
	"io"
	"math"
	"os"
)

var bitsNum = 9

func Compress(fileToCompress *os.File, compressedFilePath string, dictSize int) error {

	dict := make(map[string]int64, dictSize)
	for code := 0; code < 256; code++ {
		dict[string(rune(code))] = int64(code)
	}

	bits := int(math.Sqrt(float64(dictSize)) + 1)

	bitReader := bitstream.NewReader(fileToCompress)

	compressedFile, err := os.Create(compressedFilePath)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	bitWriter := bitstream.NewWriter(compressedFile)

	var nextCode = int64(len(dict))
	var phrase []byte
	readByte, err := bitReader.ReadByte()
	if err != nil && err != io.EOF {
		return err
	}

	phrase = append(phrase, readByte)

	for {
		char, err := bitReader.ReadByte()
		if err == io.EOF {
			finalCode := dict[string(phrase)]

			if finalCode != 0 {
				errw := bitWriter.WriteBits(uint64(finalCode), bits)
				if errw != nil && err != io.EOF {
					return errw
				}
				errw = bitWriter.Flush(false)
				if errw != nil && err != io.EOF {
					return err
				}
			}
			break
		}
		if err != nil && err != io.EOF {
			return err
		}

		withChar := string(phrase) + string(char)
		_, ok := dict[withChar]

		if !ok {
			code := dict[string(phrase)]

			err := bitWriter.WriteBits(uint64(code), bits)
			if err != nil && err != io.EOF {
				return err
			}

			if nextCode > 1<<bits {
				bits++
			}

			dict[withChar] = nextCode
			nextCode++

			phrase = nil
		}
		phrase = append(phrase, char)
	}
	return error(nil)
}

func Decompress(fileToDecompress *os.File, decompressedFilePath string, dictSize int) error {

	bits := int(math.Sqrt(float64(dictSize)) + 1)

	bitReader := bitstream.NewReader(fileToDecompress)

	decompressedFile, err := os.Create(decompressedFilePath)
	if err != nil {
		return err
	}
	defer decompressedFile.Close()

	bitWriter := bitstream.NewWriter(decompressedFile)

	dict := make(map[int64]string, dictSize)
	for code := 0; code < 256; code++ {
		dict[int64(code)] = string(rune(code))
	}

	var nextCode = int64(len(dict))

	lastCodeUns, err := bitReader.ReadBits(bits)
	if err != nil && err != io.EOF {
		return err
	}
	var lastCode = int64(lastCodeUns)

	// If file is empty.
	if lastCode == 0 {
		return error(nil)
	}

	werr := bitWriter.WriteByte(dict[lastCode][0])
	if werr != nil && werr != io.EOF {
		return werr
	}
	oldCode := lastCode

	for {

		if nextCode > 1<<bits {
			bits++
		}

		codeUnsigned, err := bitReader.ReadBits(bits)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		var code = int64(codeUnsigned)

		var outputString string
		_, ok := dict[code]

		if !ok {
			lastString := dict[oldCode]
			outputString = lastString + string(lastString[0])
		} else {
			outputString = dict[code]
		}

		for _, val := range []byte(outputString) {
			werr := bitWriter.WriteByte(val)
			if werr != nil && werr != io.EOF {
				return werr
			}
		}

		nextStringToAdd := dict[oldCode] + string(outputString[0])
		dict[nextCode] = nextStringToAdd
		nextCode++

		oldCode = code
	}
	return error(nil)
}
