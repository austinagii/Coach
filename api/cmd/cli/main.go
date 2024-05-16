package main

import (
	"aisu.ai/api/v2/cmd/cli/api"
	"aisu.ai/api/v2/cmd/cli/utils"
	"bufio"
	// "fmt"
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
		log.Printf("An error occurred while creating a new user: %s", err.Error())
		os.Exit(1)
	}
	chatId, initialMsg, err := api.CreateChat(userId)
	if err != nil {
		log.Printf("An error occurred while creating a new user: %s", err.Error())
		os.Exit(1)
	}
	utils.PromptUser(initialMsg)

	for true {
		userMessage, err := utils.ReadUserInput()
		if err != nil {
			log.Printf("An error occurred while reading the user input: %s", err.Error())
			os.Exit(1)
		}

		assistantMessage, err := api.Respond(userId, chatId, userMessage)
		if err != nil {
			log.Printf("An error occurred while reading the assistant response: %s", err.Error())
			os.Exit(1)
		}
		utils.PromptUser(assistantMessage)
	}
}
