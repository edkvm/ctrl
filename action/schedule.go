package action

import "time"

type ScheduleID string

type Schedule struct {
	ID        ScheduleID
	Action    string
	Start     time.Time
	Recurring bool
	Interval  int
}

type ScheduleRepo interface {
	Store(t *Schedule) error
	FindNext(dur time.Duration) *Schedule
	FindAllByTime(time time.Time) []*Schedule
}

func NewSchedule(action string, start time.Time) *Schedule {

	return &Schedule{
		ID:     ScheduleID(genULID()),
		Action: action,
		Start:  start,
	}
}
