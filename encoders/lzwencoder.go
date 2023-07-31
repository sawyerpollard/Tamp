package encoders

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type LZWEncoder struct {
	Alphabet             []rune  // Valid codepoints in encoding
	InitialBits          int     // Initial codeword length
	MaxBits              int     // Max codeword length
	CompressionThreshold float64 // Minimum compression ratio before flushing dictionary
}

const (
	ResizeSignal = 0
	FlushSignal  = 1
)

// LZW returns a LZWEncoder struct with an Alphabet slice and Codewords map.
func LZW(alphabet []rune, maxBits int, compressionThreshold float64) LZWEncoder {
	return LZWEncoder{
		Alphabet:             alphabet,
		InitialBits:          8,
		MaxBits:              maxBits,
		CompressionThreshold: compressionThreshold,
	}
}

// Encode implements Encoder for LZW.
func (e LZWEncoder) Encode(source Source) Code {
	pointer := 2
	dictionary := make(map[string]int)
	for _, character := range e.Alphabet {
		dictionary[string(character)] = pointer
		pointer++
	}

	uncompressedSize := 1.0
	compressedSize := 1.0

	codelength := e.InitialBits

	var sb strings.Builder
	characters := []byte(cleanSource(source, e.Alphabet))
	pattern := ""
	for i := 0; i < len(characters); i++ {
		pattern += string(characters[i])
		if dictionary[pattern] == 0 { // Pattern not in dictionary
			//fmt.Println(e.MaxBits)
			//fmt.Println(int(math.Pow(2, float64(codelength))))
			//fmt.Println(pointer)
			if codelength < e.MaxBits && pointer >= int(math.Pow(2, float64(codelength)))-2 {
				sb.WriteString(intToCodeword(ResizeSignal, codelength))
				fmt.Println(intToCodeword(ResizeSignal, codelength))
				codelength++
			}

			if pointer < int(math.Pow(2, float64(codelength)))-2 {
				dictionary[pattern] = pointer
				pointer++
			}

			entry := pattern[0 : len(pattern)-1]
			index := dictionary[entry]

			codeword := intToCodeword(index, codelength)

			sb.WriteString(codeword)

			//fmt.Printf("%v: %v (%v)\n", entry, index, codelength)

			if pointer >= int(math.Pow(2, float64(e.MaxBits)))-2 && false {
				uncompressedSize += 8 * float64(len(entry))
				compressedSize += float64(codelength)
				compressionRatio := uncompressedSize / compressedSize

				if compressionRatio < e.CompressionThreshold {
					sb.WriteString(intToCodeword(FlushSignal, codelength))
					pointer = 2
					dictionary = make(map[string]int)
					for _, character := range e.Alphabet {
						dictionary[string(character)] = pointer
						pointer++
					}
				}

				codelength = e.InitialBits
			}

			pattern = ""
			i = i - 1
		}
	}

	return Code(sb.String())
}

// Decode implements Encoder for LZW.
func (e LZWEncoder) Decode(code Code) Source {
	pointer := 2
	dictionary := make(map[string]int)
	reverseDictionary := make(map[int]string)
	for _, character := range e.Alphabet {
		dictionary[string(character)] = pointer
		reverseDictionary[pointer] = string(character)
		pointer++
	}

	codelength := e.InitialBits

	var sb strings.Builder
	pattern := ""
	for i := 0; i <= len(code)-codelength; i += codelength {
		codeword := string(code[i : i+codelength])
		parsed, _ := strconv.ParseInt(codeword, 2, 64)
		index := int(parsed)

		if index == ResizeSignal {
			codelength++
			i--
			fmt.Println("RESIZE")
			continue
		} else if index == FlushSignal {
			pointer = 2
			dictionary = make(map[string]int)
			reverseDictionary = make(map[int]string)
			for _, character := range e.Alphabet {
				dictionary[string(character)] = pointer
				reverseDictionary[pointer] = string(character)
				pointer++
			}

			codelength = e.InitialBits
			fmt.Println("FLUSH")
			continue
		}

		entry, found := reverseDictionary[index]
		if found { // Index in dictionary
			sb.WriteString(entry)

			for _, character := range entry {
				pattern += string(character)

				index, found = dictionary[pattern]
				if !found { // Pattern not in dictionary
					dictionary[pattern] = pointer
					reverseDictionary[pointer] = pattern
					pointer++

					pattern = string(character)
				}
			}
		} else {
			fmt.Println("WILD")
			fmt.Println(pointer)
			fmt.Println(index)

			for _, character := range pattern {
				pattern += string(character)

				index, found = dictionary[pattern]
				if !found { // Pattern not in dictionary
					dictionary[pattern] = pointer
					reverseDictionary[pointer] = pattern
					pointer++

					sb.WriteString(pattern)
					pattern = string(character)
				}
			}
		}
	}

	return Source(sb.String())
}

func intToCodeword(n int, length int) string {
	// Convert int to binary string
	codeword := strconv.FormatInt(int64(n), 2)

	// Pad binary string to length with 0s
	if len(codeword) < length {
		codeword = strings.Repeat("0", length-len(codeword)) + codeword
	}

	return codeword
}

func cleanSource(source Source, alphabet []rune) Source {
	valid := make(map[rune]bool)
	for _, character := range alphabet {
		valid[character] = true
	}

	var sb strings.Builder
	for _, character := range source {
		if valid[character] {
			sb.WriteString(string(character))
		}
	}

	return Source(sb.String())
}
