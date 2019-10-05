package trigger

import "time"

type ScheduleID string

type Schedule struct {
	ID     ScheduleID
	Action string
	Time time.Time
}

