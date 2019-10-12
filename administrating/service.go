package administrating

import (
	"github.com/edkvm/ctrl/action"
	"time"
)

type Service interface {
	ScheduleAction(name string, start time.Time) (action.ScheduleID, error)
	ScheduleRecurringAction(name string, interval int) (action.ScheduleID, error)
}

type service struct {
	actionRepo action.ActionRepo
	schedRepo action.ScheduleRepo
}

func NewService(actRepo action.ActionRepo, schedRepo action.ScheduleRepo) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo: schedRepo,
	}
}

func (s *service) ScheduleAction(name string, start time.Time) (action.ScheduleID, error) {
	panic("implement me")
}

func (s *service) ScheduleRecurringAction(name string, interval int) (action.ScheduleID, error) {
	panic("implement me")
}



