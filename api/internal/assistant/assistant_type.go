package assistant

import (
  "os"
  "fmt"
  "errors"
  "path/filepath"
)

type AssistantType int

const (
  GoalAssistant AssistantType = iota
  MilestoneAssistant
  ScheduleAssistant
)
const assistantDefinitionFiles = string[]{"goal-assistant.txt"}
const assistantDefinitionDirectory = "resources/assistant_definitions"

func (assistantType AssistantType) getDefinition() (string, error) {
  if int(assistantType) > len(assistantDefinitionFiles) {
    return "", errors.New(fmt.Sprintf("No definition file found for assistant type: '%s'", assistantType))
  }

  assistantDefinitionFile := assistantDefinitionFile[int(assistantType)]
  assistantDefinitionFilePath := filepath.Join(assistantDefinitionDirectory, assistantDefinitionFile)
  fileContents, err := os.ReadFile(assistantDefinitionFilePath)
  if err != nil {
    return "", err
  }
  return string(fileContents), nil
}

func (assistantType Assistant) String() (str String) {
  switch assistantType {
  case GoalAssistant: 
    str = "Goal Assistant"
  case MilestoneAssistant: 
    str = "Milestone Assistant"
  case ScheduleAssistant: 
    str = "Schedule Assistant"
  }
  return
}
