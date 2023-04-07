package main

import (
	"bytes"
	"fmt"
	"github.com/sawyerpollard/tamp/encoders"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	alphabet := Alphabet()
	corpus := Corpus("corpus")

	digramEncoder := encoders.Digram(alphabet, corpus, 256)
	EncodeFile("testing/test.txt", "digram", digramEncoder)

	huffmanEncoder := encoders.Huffman(alphabet, corpus)
	EncodeFile("testing/test.txt", "huffman", huffmanEncoder)
}

func EncodeFile(name string, extension string, e encoders.Encoder) {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	source := encoders.Source(data)

	// Print source
	fmt.Println("SOURCE TEXT:")
	fmt.Println(source, "\n")

	// Time compression
	start := time.Now()
	compressed := e.Encode(source)
	compressionTime := float32(time.Since(start)) / 10e6

	// Write compressed file
	err = WriteCode(name+"."+extension, compressed)
	if err != nil {
		log.Fatal(err)
	}

	// Read compressed file
	code, err := ReadCode(name + "." + extension)
	if err != nil {
		log.Fatal(err)
	}

	// Time decompression
	start = time.Now()
	decompressed := e.Decode(code)
	decompressionTime := float32(time.Since(start)) / 10e6

	// Print decompressed file
	fmt.Println("DECOMPRESSED TEXT:")
	fmt.Println(decompressed, "\n")

	fmt.Println("STATS:")
	fmt.Printf("COMPRESSION TIME = %v ms\n", compressionTime)
	fmt.Printf("DECOMPRESSION TIME = %v ms\n", decompressionTime)

	// Compute compression ratio
	uncompressedSize := float32(8 * len(source))
	compressedSize := float32(len(compressed))
	compressionRatio := uncompressedSize / compressedSize
	fmt.Printf("Compression ratio: %v\n", compressionRatio)
}

// WriteCode writes an encoded string to a file.
func WriteCode(name string, code encoders.Code) error {
	buffer := bytes.Buffer{}

	// Pad binary string length to multiple of 8 with 0s
	step := 8
	if len(code)%step > 0 {
		code += encoders.Code(strings.Repeat("0", step-len(code)%step))
	}

	for i := 0; i <= len(code)-step; i += step {
		// Convert binary string to byte
		parsed, err := strconv.ParseInt(string(code)[i:i+step], 2, 64)
		if err != nil {
			return err
		}

		buffer.WriteByte(byte(parsed))
	}

	err := os.WriteFile(name, buffer.Bytes(), 0644)
	return err
}

// ReadCode returns an encoded string from a file.
func ReadCode(name string) (encoders.Code, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	for _, b := range data {
		// Convert byte to binary string
		code := strconv.FormatInt(int64(b), 2)

		// Pad binary string to length 8 with 0s
		step := 8
		if len(code)%step > 0 {
			code = strings.Repeat("0", step-len(code)%step) + code
		}

		sb.WriteString(code)
	}

	return encoders.Code(sb.String()), nil
}

// Alphabet returns a slice of printable Unicode codepoints.
func Alphabet() []rune {
	alphabet := make([]rune, 0, 128)

	alphabet = append(alphabet, rune(9))  // Tab character
	alphabet = append(alphabet, rune(10)) // Newline character
	for i := 32; i <= 126; i++ {
		alphabet = append(alphabet, rune(i)) // Printable characters
	}

	return alphabet
}

// Corpus returns a concatenated string of text files from directory.
func Corpus(directory string) string {
	corpus := ""

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(path.Join(directory, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		corpus += string(data)
	}

	return corpus
}
