package main

import (
	"aisu.ai/api/v2/cmd/cli/api"
	"aisu.ai/api/v2/cmd/cli/utils"
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	inputReader = bufio.NewReader(os.Stdin)
)

func main() {
	utils.PromptUser("What's your name? Dont play with me either!")
	name, err := utils.ReadUserInput()
	if err != nil {
		// log.Printf("An error occurred while reading user input: %s\n", err.Error())
		os.Exit(1)
	}
	userId, err := api.CreateUser(name)
	if err != nil {
		log.Printf("An error occurred while creating a new user: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("User created with ID '%s'\n", userId)
}
