package task

import (
	"fmt"
	"log/slog"
)

type Task struct {
	Objective  Objective `json:"objective" bson:"objective"`
	IsComplete bool      `json:"is_complete" bson:"is_complete"`
	TargetId   string    `json:"target_id,omitempty"`
}

// NewTask returns a new [Task] with the specified objective.
// The target ID is optional and should only be specified if the
// goal is milestone creation and in this case the target ID should
// be the ID of the goal whose milestones are to be defined.
func NewTask(objective Objective, targetId string) *Task {
	return &Task{
		Objective:  objective,
		IsComplete: false,
		TargetId:   targetId,
	}
}

// Description returns the description of the task, dynamically formatted based on the objective.
func (t Task) Description() (string, error) {
	description, err := t.Objective.Description()
	if err != nil {
		const errMsg = "Failed to retrieve objective description"
		slog.Error(errMsg, "error", err)
		return "", fmt.Errorf("%s: %w", errMsg, err)
	}

	// Format the description of the milestone creation objective
	// to include the ID of the goal for which the milestones are
	// being created
	if t.Objective == ObjectiveMilestoneCreation {
		description = fmt.Sprintf(description, t.TargetId)
	}

	return description, nil
}
