package administrating

import (
	"fmt"
	"github.com/edkvm/ctrl/action"
	"time"
)

type Service interface {
	ListSchedule(name string) ([]*action.Schedule, error)
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

func (s *service) ListSchedule(name string) ([]*action.Schedule, error) {

	return s.schedRepo.FindAllByAction(name), nil
}

func (s *service) ScheduleAction(name string, start time.Time) (action.ScheduleID, error) {
	panic("implement me")
}

func (s *service) ScheduleRecurringAction(name string, interval int) (action.ScheduleID, error) {
	sched := action.NewRecurringSchedule(name, time.Now(), interval)
	err := s.schedRepo.Store(sched)
	if err != nil {
		return "", fmt.Errorf("could not save")
	}

	return sched.ID, nil
}



