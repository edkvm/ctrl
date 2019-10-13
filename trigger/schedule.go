package trigger

import (

	"time"

	"github.com/edkvm/ctrl/action"
	ctrlID "github.com/edkvm/ctrl/pkg/id"
)

type ScheduleID string

type Schedule struct {
	ID        ScheduleID
	Action    string
	Start     time.Time
	Recurring bool
	Interval  int
	Params    action.ActionParams
}

type SchedulingEvent struct {
	ScheduleID ScheduleID
	Action string
}

type ScheduleRepo interface {
	Store(t *Schedule) error
	FindNext(dur time.Duration) *Schedule
	Find(id ScheduleID) (*Schedule, error)
	FindAllByTime(time time.Time) []*Schedule
	FindAllByAction(name string) []*Schedule
}

func NewSchedule(action string, start time.Time) *Schedule {
	return &Schedule{
		ID:     ScheduleID(ctrlID.GenULID()),
		Action: action,
		Start:  start,
		Recurring: false,
	}
}

func NewRecurringSchedule(action string, start time.Time, interval int, params action.ActionParams) *Schedule {
	return &Schedule{
		ID:     ScheduleID(ctrlID.GenULID()),
		Action: action,
		Start:  start,
		Recurring: true,
		Interval: interval,
		Params: params,
	}
}


