package invoking

import "github.com/edkvm/ctrl/action"

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	ScheduleAction(name string, params map[string]interface{}) (interface{}, error)
}

type service struct {
	actionRepo action.ActionRepo
	schedRepo action.ScheduleRepo
	provider *action.ActionProvider
}

func NewService(actRepo action.ActionRepo, schedRepo action.ScheduleRepo, provider *action.ActionProvider) *service {
	return &service{
		actionRepo: actRepo,
		schedRepo: schedRepo,
		provider: provider,
	}
}

func (s *service) ScheduleAction(name string, params map[string]interface{}) (interface{}, error) {
	panic("implement me")
}

func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {


	ac := action.NewAction(name)

	env := ac.BuildEnv()
	payload := s.provider.EncodePayload(params)
	result := s.provider.ExecuteAction(name, payload, env)

	return result, nil
}

