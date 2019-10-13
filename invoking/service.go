package invoking

import (
	"github.com/edkvm/ctrl/action"
	"time"
)

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	AddActionSchedule(name string, schedID action.ScheduleID) error
}

type service struct {
	actionRepo action.ActionRepo
	schedRepo action.ScheduleRepo
	statsRepo action.StatsRepo
	actionTimer *action.ActionTimer
	provider *action.ActionProvider
}

func NewService(actRepo action.ActionRepo, schedRepo action.ScheduleRepo, statsRepo action.StatsRepo, actionTimer *action.ActionTimer, provider *action.ActionProvider) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo: schedRepo,
		statsRepo: statsRepo,
		actionTimer: actionTimer,
		provider: provider,
	}
}

func (s *service) AddActionSchedule(name string, schedID action.ScheduleID) error {
	sched, _ := s.schedRepo.Find(schedID)
	s.actionTimer.AddSchedule(sched, s.RunAction)

	return nil
}

func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {
	defer func(name string, start time.Time){

		stat := action.NewStat(name, start, time.Now(), true)
		s.statsRepo.Store(stat)
	}(name, time.Now())

	ac := action.NewAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.ExecuteAction(name, payload, env)

	return result, nil
}



