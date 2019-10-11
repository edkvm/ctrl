package action

import "time"

type ScheduleID string

type ScheduleRepo interface {
	Store(t *Schedule) error
	FindNext(dur time.Duration) *Schedule
	FindAllByTime(time time.Time) []*Schedule
}


type Schedule struct {
	ID     ScheduleID
	Action string
	When   time.Time
}

func NewSchedule (action string, schedTime time.Time) *Schedule{

	return &Schedule{
		ID:     ScheduleID(genULID()),
		Action: action,
		When:   schedTime,
	}
}
