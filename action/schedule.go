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
	FindAllByAction(name string) []*Schedule
	FindAllByTime(time time.Time) []*Schedule
}

func NewSchedule(action string, start time.Time) *Schedule {
	return &Schedule{
		ID:     ScheduleID(genULID()),
		Action: action,
		Start:  start,
		Recurring: false,
	}
}

func NewRecurringSchedule(action string, start time.Time, interval int) *Schedule {
	return &Schedule{
		ID:     ScheduleID(genULID()),
		Action: action,
		Start:  start,
		Recurring: true,
		Interval: interval,
	}
}


