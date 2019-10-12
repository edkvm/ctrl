package administrating

import (
	"context"
	"github.com/edkvm/ctrl/pkg/endpoint"
	"time"
)

type scheduleRequest struct {
	Action    string
	Start     time.Time
	Recurring bool
	Interval  int
}

func makeScheduleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		schedReq := req.(scheduleRequest)
		result, err := s.ScheduleRecurringAction(schedReq.Action, schedReq.Interval)
		if err != nil {

		}
		return result, nil
	}
}