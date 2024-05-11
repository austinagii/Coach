package task

import (
	"aisu.ai/api/v2/internal/user"
	"log/slog"
)

type Task interface {
	Objective() Objective
	Description() string
}

type BaseTask struct {
	Obj Objective `json:"objective" bson:"objective"`
}

func (task *BaseTask) Objective() Objective {
	return task.Obj
}

// Description returns the description of the task, dynamically formatted based on the objective.
func (task *BaseTask) Description() string {
	description, err := task.Obj.description()
	if err != nil {
		slog.Error("Failed to retrieve objective description", "error", err)
		return ""
	}
	return description
}

type GoalCreationTask struct {
	BaseTask `bson:",inline"`
	Goal     *user.Goal `json:"goal,omitempty" bson:"goal,omitempty"`
}

func NewGoalCreationTask() *GoalCreationTask {
	return &GoalCreationTask{BaseTask: BaseTask{Obj: ObjectiveGoalCreation}}
}

type MilestoneCreationTask struct {
	BaseTask   `bson:",inline"`
	GoalId     int               `json:"goal_id" bson:"goal_id"`
	Milestones []*user.Milestone `json:"milestones,omitempty" bson:"milestones,omitempty"`
}

func NewMilestoneCreationTask(goalId int) *MilestoneCreationTask {
	return &MilestoneCreationTask{BaseTask: BaseTask{Obj: ObjectiveMilestoneCreation}, GoalId: goalId}
}

type ScheduleCreationTask struct {
	BaseTask `bson:",inline"`
	Schedule user.DailySchedule `json:"schedule,omitempty" bson:"schedule,omitempty"`
}

func NewScheduleCreationTask() *ScheduleCreationTask {
	return &ScheduleCreationTask{BaseTask: BaseTask{Obj: ObjectiveScheduleCreation}}
}
