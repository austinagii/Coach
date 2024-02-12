package main

import (
  "fmt"
  "net/http"
  // "io/ioutil"
  aisu "aisu.ai/api/v2/internal/assistant"
  "github.com/gin-gonic/gin"
)

type GoalCompletionRequest struct {
  UserMessage string `json:"user_message"`
}

type GoalCompletionResponse struct {
  AssistantMessage  string  `json:"assistant_message"`
  Goal              string  `json:"goal"`
  IsGoalDefined     bool    `json:"is_goal_defined"`
}

var assistant *aisu.Assistant = aisu.NewAssistant(`You are ANa, an assistant created to help users define a single SMART goal. 
    You are to prompt the user for their goal and guide them through making it a SMART goal if it 
    is not one already. SMART in this case being the acronym for Specific Measurable Achievable, 
    Relevant and Time-Bound. On every response, you should respond in RAW JSON (do not include markdown 
    tags) with this structure: {"goal": "<the users goal>", "is_smart_goal": <boolean indicating whether 
    the goal is a smart goal>, "message": {"text": "<the next message to help the user make their goal SMART>"}}. 
    Assume the user doesnt know how to break down their goal to make it SMART, you are to ask the user 
    questions about their goal to help them turn their goal into a SMART goal. Do not explicitly mention 
    SMART in your responses.`)

func corsMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

    // Handle preflight requests
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(http.StatusOK)
        return
    }

    c.Next()
  }
}

func getGoalCompletion(context *gin.Context) {
  // bodyBytes, _ := ioutil.ReadAll(context.Request.Body)
  // rawJSON := string(bodyBytes)
  // fmt.Println(rawJSON)
  // return 

  var request GoalCompletionRequest

  if err := context.BindJSON(&request); err != nil {
    panic(err)
  }

  userMessage := aisu.NewUserMessage(request.UserMessage)
  fmt.Println(userMessage)
  assistantResponse, err := assistant.Respond(userMessage)
  if err != nil {
    panic(err)
  }

  response := GoalCompletionResponse{
    AssistantMessage: assistantResponse.Message.Text,
    Goal: assistantResponse.Goal,
    IsGoalDefined: assistantResponse.IsSmartGoal,
  }

  context.IndentedJSON(http.StatusOK, response)
}

func main() {
    router := gin.Default()
    router.Use(corsMiddleware())
    router.POST("/goal-completion", getGoalCompletion)
    router.Run("0.0.0.0:8080")
}
