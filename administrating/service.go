package administrating

import (
	"fmt"
	"github.com/edkvm/ctrl/action"
	"github.com/edkvm/ctrl/invoking"
	"time"
)

type EventHandler interface {
	ActionWasScheduled(event action.SchedulingEvent)
}

type Service interface {
	ListSchedule(name string) ([]*action.Schedule, error)
	ScheduleAction(name string, start time.Time) (action.ScheduleID, error)
	ScheduleRecurringAction(name string, interval int, params action.ActionParams) (action.ScheduleID, error)

	ListStats(name string) ([]*action.Stat, error)
}

type service struct {
	actionRepo action.ActionRepo
	schedRepo  action.ScheduleRepo
	statsRepo  action.StatsRepo
	schedHandler EventHandler
}

func NewService(actRepo action.ActionRepo, schedRepo action.ScheduleRepo, statsRepo action.StatsRepo,schedHandler EventHandler) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo:  schedRepo,
		statsRepo: statsRepo,
		schedHandler: schedHandler,
	}
}

func (s *service) ListSchedule(name string) ([]*action.Schedule, error) {

	return s.schedRepo.FindAllByAction(name), nil
}

func (s *service) ListStats(name string) ([]*action.Stat, error) {

	return s.statsRepo.FindByAction(name), nil
}

func (s *service) ScheduleAction(name string, start time.Time) (action.ScheduleID, error) {
	panic("implement me")
}

func (s *service) ScheduleRecurringAction(name string, interval int, params action.ActionParams) (action.ScheduleID, error) {
	sched := action.NewRecurringSchedule(name, time.Now(), interval, params)
	err := s.schedRepo.Store(sched)
	if err != nil {
		return "", fmt.Errorf("could not save")
	}

	s.schedHandler.ActionWasScheduled(action.SchedulingEvent{
		Action: name,
		ScheduleID: sched.ID,
	})

	return sched.ID, nil
}

type schedulingEventHandler struct {
	InvokingService invoking.Service
}

// TODO: Reduce dependency with other service by using the func only without the service and
//       create the func in main with options to use http, for distribution.
func (h *schedulingEventHandler) ActionWasScheduled(event action.SchedulingEvent) {
	h.InvokingService.AddActionSchedule(event.Action, event.ScheduleID)
}

func NewEventHandler(s invoking.Service) EventHandler {
	return &schedulingEventHandler{
		InvokingService: s,
	}
}

