package task

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Objective string

const (
	ObjectiveGoalCreation      Objective = "goal_creation"
	ObjectiveMilestoneCreation Objective = "milestone_creation"
	ObjectiveScheduleCreation  Objective = "schedule_creation"
	ObjectiveChat              Objective = "chat"
)

var availableObjectives = []Objective{
	ObjectiveGoalCreation,
	ObjectiveMilestoneCreation,
	ObjectiveScheduleCreation,
	ObjectiveChat,
}
var descriptionByObjective = map[Objective]string{}

func LoadObjectiveDescriptions() error {
	descriptionFileByObjective := map[Objective]string{
		ObjectiveGoalCreation:      "../../resources/assistant/objectives/goal_creation/description.txt",
		ObjectiveMilestoneCreation: "../../resources/assistant/objectives/milestone_creation/description.txt",
		ObjectiveScheduleCreation:  "../../resources/assistant/objectives/schedule_creation/description.txt",
		ObjectiveChat:              "../../resources/assistant/objectives/chat/description.txt",
	}

	for objective, descriptionFile := range descriptionFileByObjective {
		fileContents, err := os.ReadFile(descriptionFile)
		if err != nil {
			errMsg := "An error occurred while loading the objective's description file"
			slog.Error(errMsg, "objective", objective, "error", err)
			return fmt.Errorf("%s: %w", errMsg, err)
		}
		descriptionByObjective[objective] = string(fileContents)
	}
	return nil
}

func (o Objective) Description() (string, error) {
	description, ok := descriptionByObjective[o]
	if !ok {
		return "", fmt.Errorf("No description found for objective '%s'", o)
	}
	return description, nil
}

func (o Objective) String() (str string) {
	return string(o)
}

func objectiveFromString(s string) (Objective, error) {
	for _, objective := range availableObjectives {
		if objective.String() == s {
			return objective, nil
		}
	}

	return "", fmt.Errorf("Invalid objective: '%s'. Valid objectives are: %s", s, allAvailableObjectiveStrings())
}

func (o Objective) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *Objective) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	objective, err := objectiveFromString(s)
	if err != nil {
		return err
	}
	*o = objective
	return nil
}

func allAvailableObjectiveStrings() string {
	objectives := make([]string, 0, len(availableObjectives))
	for _, obj := range availableObjectives {
		objectives = append(objectives, fmt.Sprintf("'%s'", obj))
	}
	return strings.Join(objectives, ", ")
}
