package invoking

import (
	"github.com/edkvm/ctrl/action"
	"github.com/edkvm/ctrl/execute"
	"github.com/edkvm/ctrl/trigger"
	"time"
)

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	AddActionSchedule(name string, schedID trigger.ScheduleID) error
}

type service struct {
	actionRepo  action.ActionRepo
	schedRepo   trigger.ScheduleRepo
	statsRepo   action.StatsRepo
	actionTimer *execute.ActionTimer
	provider    *execute.ActionProvider
}

func NewService(actRepo action.ActionRepo, schedRepo trigger.ScheduleRepo, statsRepo action.StatsRepo, actionTimer *execute.ActionTimer, provider *execute.ActionProvider) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo: schedRepo,
		statsRepo: statsRepo,
		actionTimer: actionTimer,
		provider: provider,
	}
}

func (s *service) AddActionSchedule(name string, schedID trigger.ScheduleID) error {
	sched, _ := s.schedRepo.Find(schedID)
	s.actionTimer.AddSchedule(sched, s.RunAction)

	return nil
}

func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {
	stat := action.NewStat(name, time.Now(), action.Running)
	defer func(stat *action.Stat){
		s.statsRepo.Store(stat)
	}(stat)

	// TODO: Add Error handeling
	ac := action.NewAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.ExecuteAction(name, payload, env)

	return result, nil
}



