package action

import "time"

type ScheduleID string

type Schedule struct {
	ID     ScheduleID
	Action string
	When   time.Time
}

func NewSchedule (action string, schedTime time.Time) *Schedule{

	return &Schedule{
		ID:     ScheduleID(ctrl.GenULID()),
		Action: action,
		When:   schedTime,
	}
}
