package main

import (
  aisu "aisu.ai/api/v2/internal/assistant"
  "bufio"
  "os"
  "fmt"
)

func main() {
  assistant := aisu.NewAssistant(`You are ANa, an assistant created to help users define a single SMART goal. 
    You are to prompt the user for their goal and guide them through making it a SMART goal if it 
    is not one already. SMART in this case being the acronym for Specific Measurable Achievable, 
    Relevant and Time-Bound. On every response, you should respond in RAW JSON (do not include markdown 
    tags) with this structure: {"goal": "<the users goal>", "is_smart_goal": <boolean indicating whether 
    the goal is a smart goal>, "message": {"text": "<the next message to help the user make their goal SMART>"}}. 
    Assume the user doesnt know how to break down their goal to make it SMART, you are to ask the user 
    questions about their goal to help them turn their goal into a SMART goal. Do not explicitly mention 
    SMART in your responses.`)

  // Print all messages except the system message
  for iterator := assistant.Chat.Front().Next(); iterator != nil; iterator = iterator.Next() {
    fmt.Println(iterator.Value)
  }

  // Read the user's input
  inputReader := bufio.NewReader(os.Stdin)
  for {
    userInput, err := readUserInput(inputReader)
    if err != nil {
      panic(err)
    }

    if userInput == "bye" || userInput == "exit" {
      break
    } 

    assistantResponse, err := assistant.Respond(aisu.NewUserMessage(userInput))
    if err != nil {
      panic(err)
    }

    if assistantResponse.IsSmartGoal {
      fmt.Printf("\nGreat! Your goal is: %v\n", assistantResponse.Goal)
      break
      // return assistantResponse.Goal
    } else {
      fmt.Println(assistantResponse.Message)
    }
  }
}

func readUserInput(reader *bufio.Reader) (string, error) {
  fmt.Printf("User: ")
  input, err := reader.ReadString('\n')
  input = input[:len(input) - 1]
  if err != nil {
    return "", nil
  }
  return input, nil
}
