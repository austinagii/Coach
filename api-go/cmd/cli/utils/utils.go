package utils

import (
	"bufio"
	"fmt"
	"os"
)

var inputReader = bufio.NewReader(os.Stdin)

func ReadUserInput() (string, error) {
	fmt.Printf("User: ")
	input, err := inputReader.ReadString('\n')
	input = input[:len(input)-1]
	if err != nil {
		return "", nil
	}
	return input, nil
}

func PromptUser(prompt string) {
	fmt.Printf("AISU: %s\n", prompt)
}
