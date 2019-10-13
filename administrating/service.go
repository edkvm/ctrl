package administrating

import (
	"fmt"
	"time"

	"github.com/edkvm/ctrl/action"
	"github.com/edkvm/ctrl/trigger"
)

type EventHandler interface {
	ActionWasScheduled(event trigger.SchedulingEvent)
}

type Service interface {
	ListSchedule(name string) ([]*trigger.Schedule, error)
	ScheduleAction(name string, start time.Time) (trigger.ScheduleID, error)
	ScheduleRecurringAction(name string, interval int, params action.ActionParams) (trigger.ScheduleID, error)
	ToggleSchedule(id trigger.ScheduleID) error

	ListStats(name string) ([]*action.Stat, error)
}

type service struct {
	actionRepo   action.ActionRepo
	schedRepo    trigger.ScheduleRepo
	statsRepo    action.StatsRepo
	schedHandler EventHandler
}

func NewService(actRepo action.ActionRepo, schedRepo trigger.ScheduleRepo, statsRepo action.StatsRepo,schedHandler EventHandler) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo:  schedRepo,
		statsRepo: statsRepo,
		schedHandler: schedHandler,
	}
}

func (s *service) ListSchedule(name string) ([]*trigger.Schedule, error) {

	return s.schedRepo.FindAllByAction(name), nil
}

func (s *service) ListStats(name string) ([]*action.Stat, error) {

	return s.statsRepo.FindByAction(name), nil
}

func (s *service) ScheduleAction(name string, start time.Time) (trigger.ScheduleID, error) {
	panic("implement me")
}

func (s *service) ScheduleRecurringAction(name string, interval int, params action.ActionParams) (trigger.ScheduleID, error) {
	sched := trigger.NewRecurringSchedule(name, time.Now(), interval, params)
	err := s.schedRepo.Store(sched)
	if err != nil {
		return "", fmt.Errorf("could not save")
	}

	s.schedHandler.ActionWasScheduled(trigger.SchedulingEvent{
		Action: name,
		ScheduleID: sched.ID,
	})

	return sched.ID, nil
}

func (s *service) ToggleSchedule(id trigger.ScheduleID) error {
	sched, err := s.schedRepo.Find(id)
	if err != nil {
		return err
	}

	sched.Active = !sched.Active

	return nil
}

type schedulingEventHandler struct {
	HandlerFunc func (name string, schedID trigger.ScheduleID) error
}

// TODO: Reduce dependency with other service by using the func only without the service and
//       create the func in main with options to use http, for distribution.
func (h *schedulingEventHandler) ActionWasScheduled(event trigger.SchedulingEvent) {
	h.HandlerFunc(event.Action, event.ScheduleID)
}

func NewEventHandler(h func (name string, schedID trigger.ScheduleID) error) EventHandler {
	return &schedulingEventHandler{
		HandlerFunc: h,
	}
}

