package encoders

import (
	"math"
	"strconv"
	"strings"
)

type DigramEncoder struct {
	Alphabet  []rune            // Valid codepoints in encoding
	Codewords map[string]string // Map from codepoints to codewords
	Size      int               // Dictionary size
}

// Digram returns a DigramEncoder struct with an Alphabet slice and Codewords map.
func Digram(alphabet []rune, corpus string, dictionarySize int) DigramEncoder {
	alphabetSize := len(alphabet)

	frequencies := digramFrequencies(corpus, alphabet)
	digrams := mostFrequentDigrams(dictionarySize-alphabetSize, frequencies)

	var i int64 = 0
	codewords := make(map[string]string)
	for _, character := range alphabet {
		// Convert int to binary string
		codeword := strconv.FormatInt(i, 2)

		codelength := int(math.Log2(float64(dictionarySize)))

		// Pad binary string to codelength with 0s
		if len(codeword) < codelength {
			codeword = strings.Repeat("0", codelength-len(codeword)) + codeword
		}

		codewords[string(character)] = codeword
		i++
	}

	for _, digram := range digrams {
		// Convert int to binary string
		codeword := strconv.FormatInt(i, 2)

		codelength := int(math.Log2(float64(dictionarySize)))

		// Pad binary string to codelength with 0s
		if len(codeword) < codelength {
			codeword = strings.Repeat("0", codelength-len(codeword)) + codeword
		}

		codewords[digram] = codeword
		i++
	}

	return DigramEncoder{
		Alphabet:  alphabet,
		Codewords: codewords,
		Size:      dictionarySize,
	}
}

func mostFrequentDigrams(limit int, frequencies map[string]int) []string {
	frequenciesCopy := make(map[string]int)
	for digram, frequency := range frequencies {
		frequenciesCopy[digram] = frequency
	}
	frequencies = frequenciesCopy

	mostFrequent := make([]string, 0)

	for len(mostFrequent) < limit && len(frequencies) > 0 {
		maxDigram := ""
		maxFrequency := 0
		for digram, frequency := range frequencies {
			if frequency > maxFrequency {
				maxDigram = digram
				maxFrequency = frequency
			}
		}

		mostFrequent = append(mostFrequent, maxDigram)
		delete(frequencies, maxDigram)
	}

	return mostFrequent
}

func digramFrequencies(text string, alphabet []rune) map[string]int {
	valid := make(map[rune]bool)
	for _, character := range alphabet {
		valid[character] = true
	}

	frequencies := make(map[string]int)

	runes := []rune(text)
	for i := 0; i < len(runes)-1; i++ {
		if valid[runes[i]] && valid[runes[i+1]] {
			digram := string(runes[i]) + string(runes[i+1])
			frequencies[digram]++
		}
	}

	return frequencies
}

// Encode implements Encoder for digram codes.
func (e DigramEncoder) Encode(source Source) Code {
	var sb strings.Builder

	// Iterate characters of source and concatenate codewords
	characters := []byte(source)
	for i := 0; i < len(characters); i++ {
		// Special case to handle last character
		if i == len(characters)-1 {
			sb.WriteString(e.Codewords[string(characters[i])])
		} else {
			digram := string(characters[i]) + string(characters[i+1])
			if e.Codewords[digram] != "" {
				sb.WriteString(e.Codewords[digram])
				i++
			} else {
				sb.WriteString(e.Codewords[string(characters[i])])
			}
		}
	}

	return Code(sb.String())
}

// Decode implements Encoder for digram codes.
func (e DigramEncoder) Decode(code Code) Source {
	// Invert codewords map
	decodewords := make(map[string]string)
	for s, codeword := range e.Codewords {
		decodewords[codeword] = s
	}

	var sb strings.Builder

	step := int(math.Log2(float64(e.Size)))
	for i := 0; i <= len(code)-step; i += step {
		codeword := string(code[i : i+step])
		sb.WriteString(decodewords[codeword])
	}

	return Source(sb.String())
}
