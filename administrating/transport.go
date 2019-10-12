package administrating

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	ctrlhttp "github.com/edkvm/ctrl/pkg/http"
)

func MakeHandler(srv Service) http.Handler {
	r := httprouter.New()

	server := ctrlhttp.NewServer(
		makeScheduleEndpoint(srv),
		decodeScheduleRequest,
		encodeInvokeResponse,
	)

	r.Handler(http.MethodPost, "/admin/v1/schedule", server)

	return r
}

func decodeScheduleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Action    string `json:"action"`
		Start     time.Time `json:"start"`
		Recurring bool `json:"recurring"`
		Interval  int	`json:"interval"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return scheduleRequest{
		Action: body.Action,
		Start: body.Start,
		Recurring: body.Recurring,
		Interval: body.Interval,
	}, nil
}

func encodeInvokeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}


