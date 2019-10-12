package invoking

import (
	"github.com/edkvm/ctrl/action"
	"time"
)

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	ScheduleAction(name string, params map[string]interface{}) (interface{}, error)
}

type service struct {
	actionRepo action.ActionRepo
	schedRepo action.ScheduleRepo
	statsRepo action.StatsRepo
	provider *action.ActionProvider
}

func NewService(actRepo action.ActionRepo, schedRepo action.ScheduleRepo, statsRepo action.StatsRepo, provider *action.ActionProvider) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo: schedRepo,
		statsRepo: statsRepo,
		provider: provider,
	}
}

func (s *service) ScheduleAction(name string, params map[string]interface{}) (interface{}, error) {
	panic("implement me")
}

func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {
	defer func(name string, start time.Time){

		stat := &action.Stat{
			ActionName: name,
			Start: start,
			End: time.Now(),
			Status: true,
		}
		s.statsRepo.Store(stat)
	}(name, time.Now())

	ac := action.NewAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.ExecuteAction(name, payload, env)

	return result, nil
}

