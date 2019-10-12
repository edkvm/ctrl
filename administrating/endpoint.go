package administrating

import (
	"context"
	"github.com/edkvm/ctrl/pkg/endpoint"
	"log"

	"time"
)

type scheduleRequest struct {
	Action    string
	Start     time.Time
	Recurring bool
	Interval  int
}

type scheduleResponse struct {
	ID 		string      `json:"id"`
	Action    string    `json:"action,omitempty"`
	Start     time.Time `json:"start,omitempty"`
	Recurring bool      `json:"recurring"`
	Interval  int       `json:"interval,omitempty"`
}

func makeCreateScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		log.Println("dd")
		schedReq := req.(scheduleRequest)
		result, err := s.ScheduleRecurringAction(schedReq.Action, schedReq.Interval)
		if err != nil {

		}
		return result, nil
	}
}

type listSchedule struct {
	Schedules []scheduleResponse `json:"schedules,omitempty"`
}

func makeListScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		name := req.(string)

		result, err := s.ListSchedule(name)
		if err != nil {

		}

		list := make([]scheduleResponse,  0, len(result))
		for _, val := range result {
			list = append(list, scheduleResponse{
				string(val.ID),
				val.Action,
				val.Start,
				val.Recurring,
				val.Interval,
			})
		}

		return listSchedule{Schedules: list}, nil
	}
}
