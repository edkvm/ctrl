package administrating

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	ctrlhttp "github.com/edkvm/ctrl/pkg/http"
)

func MakeHandler(srv Service) http.Handler {
	r := httprouter.New()

	createScheduleHandler := ctrlhttp.NewServer(
		makeCreateScheduleEndpoint(srv),
		decodeScheduleRequest,
		encodeResponse,
	)

	listScheduleHandler := ctrlhttp.NewServer(
		makeListScheduleEndpoint(srv),
		decodeActionName,
		encodeResponse,
	)

	listStatsHandler := ctrlhttp.NewServer(
		makeListStatsEndpoint(srv),
		decodeActionName,
		encodeResponse,
	)

	r.Handler(http.MethodPost, "/admin/v1/schedule", createScheduleHandler)
	r.Handler(http.MethodGet, "/admin/v1/schedule/:name", listScheduleHandler)
	r.Handler(http.MethodGet, "/admin/v1/stats/:name", listStatsHandler)

	return r
}

func decodeScheduleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	log.Println("dec")
	var body struct {
		Action    string `json:"action"`
		Start     time.Time `json:"start"`
		Recurring bool `json:"recurring"`
		Interval  int	`json:"interval"`
		Params    map[string]interface{} `json:"params"`

	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return scheduleRequest{
		Action: body.Action,
		Start: body.Start,
		Recurring: body.Recurring,
		Interval: body.Interval,
		Params: body.Params,
	}, nil
}

func decodeActionName(ctx context.Context, r *http.Request) (interface{}, error) {
	params := httprouter.ParamsFromContext(ctx)
	name := params.ByName("name")

	return name, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	//if e, ok := response.(errorer); ok && e.error() != nil {
	//	encodeError(ctx, e.error(), w)
	//	return nil
	//}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}




