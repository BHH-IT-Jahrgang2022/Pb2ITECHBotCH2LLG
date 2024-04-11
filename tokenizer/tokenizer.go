package tokenizer

import (
	"fmt"
	"strings"
)

func main() {
	input := "Hello, world! This is a tokenizer example."

	// Tokenize the input string
	tokens := tokenize(input)

	// Print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func tokenize(input string) []string {
	// Split the input string into tokens
	tokens := strings.Fields(input)

	return tokens
}

//yolo
