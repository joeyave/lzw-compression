package lzw

import (
	"io"
	"os"

	"github.com/dgryski/go-bitstream"
)

//Compress performs compression according to parameters in Config
func Compress(fileToCompress *os.File, compressedFile *os.File, dictionary map[string]int64) error {

	bitReader := bitstream.NewReader(fileToCompress)
	bitWriter := bitstream.NewWriter(compressedFile)

	var nextCode int64 = 256
	var phrase []byte
	readByte, err := bitReader.ReadByte()
	if err != nil && err != io.EOF {
		return err
	}

	phrase = append(phrase, readByte)

	for {
		char, err := bitReader.ReadByte()
		if err == io.EOF {
			finalCode := dictionary[string(phrase)]

			if finalCode != 0 {
				errw := bitWriter.WriteBits(uint64(finalCode), conf.codeWidth)
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
		_, ok := dictionary[withChar]

		if !ok {
			code := dictionary[string(phrase)]
			err := bitWriter.WriteBits(uint64(code), conf.codeWidth)
			if err != nil && err != io.EOF {
				return err
			}

			if nextCode < conf.maxCode {
				dictionary[withChar] = nextCode
				nextCode++
			}
			phrase = nil
		}
		phrase = append(phrase, char)
	}
	return error(nil)
}

// Decompress performs decompression according to parameters in Config
func Decompress(fileToDecompress *os.File, decompressedFile *os.File, origDictionary map[string]int64) error {

	bitReader := bitstream.NewReader(fileToDecompress)
	bitWriter := bitstream.NewWriter(decompressedFile)

	dictionary := make(map[int64]string)
	for key, val := range origDictionary {
		dictionary[val] = key
	}

	var nextCode int64 = 256

	lastCodeUns, err := bitReader.ReadBits(conf.codeWidth)
	if err != nil && err != io.EOF {
		return err
	}
	var lastCode int64 = int64(lastCodeUns)

	// empty file
	if lastCode == 0 {
		return error(nil)
	}

	werr := bitWriter.WriteByte(dictionary[lastCode][0])
	if werr != nil && werr != io.EOF {
		return werr
	}
	oldCode := lastCode

	for {
		codeUnsigned, err := bitReader.ReadBits(conf.codeWidth)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		var code int64 = int64(codeUnsigned)

		var outputString string
		if _, found := dictionary[code]; !found {
			lastString := dictionary[oldCode]
			outputString = lastString + string(lastString[0])
		} else {
			outputString = dictionary[code]
		}
		for _, val := range []byte(outputString) {
			werr := bitWriter.WriteByte(val)
			if werr != nil && werr != io.EOF {
				return werr
			}
		}
		if nextCode < conf.maxCode {
			nextStringToAdd := dictionary[oldCode] + string(outputString[0])
			dictionary[nextCode] = nextStringToAdd
			nextCode++
		}
		oldCode = code
	}
	return error(nil)
}
