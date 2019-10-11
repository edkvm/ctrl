package invoking

import "github.com/edkvm/ctrl/action"

type Service interface {
	RunAction(name string, params map[string]interface{}) (interface{}, error)
	ScheduleAction(name string, params map[string]interface{}) (interface{}, error)
}

type service struct {
	actionRepo action.ActionRepo
}

func NewService(actRepo action.ActionRepo, schedRepo action.ActionRepo) *service {
	return &service{
		actionRepo: actRepo,
	}
}

func (s *service) ScheduleAction(name string, params map[string]interface{}) (interface{}, error) {
	panic("implement me")
}



func (s *service) RunAction(name string, params map[string]interface{}) (interface{}, error) {

}

