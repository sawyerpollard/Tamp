package encoders

import "strings"

type PrefixEncoder struct {
	Alphabet  []rune          // Valid codepoints in encoding
	Codewords map[rune]string // Map from codepoints to codewords
}

// Encode implements Encoder for prefix codes.
func (e PrefixEncoder) Encode(source Source) Code {
	var sb strings.Builder

	// Iterate characters of source and concatenate codewords
	for _, character := range source {
		sb.WriteString(e.Codewords[character])
	}

	return Code(sb.String())
}

// Decode implements Encoder for prefix codes.
func (e PrefixEncoder) Decode(code Code) Source {
	// Invert codewords map
	decodewords := make(map[string]rune)
	for character, codeword := range e.Codewords {
		decodewords[codeword] = character
	}

	var sb strings.Builder

	codeword := ""
	for _, character := range code {
		codeword += string(character)

		// Greedily decode binary string when it maps to a character
		if decodewords[codeword] != 0 {
			sb.WriteString(string(decodewords[codeword]))
			codeword = ""
		}
	}

	return Source(sb.String())
}
