package invoking

import (
	"context"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/edkvm/ctrl/pkg/endpoint"
	ctrlhttp "github.com/edkvm/ctrl/pkg/http"
)


type invokeRequest map[string]interface{}



func MakeHandler(srv Service) {
	r := httprouter.New()

	server := ctrlhttp.NewServer(
		makeInvokeActionEndpoint(srv),
		decodeInvokeRequest,
		encodeInvokeResponse,
	)

	r.Handler(http.MethodPost, "/action/:name", server)

}


func decodeInvokeRequest(_ context.Context, r *http.Request) (interface{}, error){
	var params map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}


	return params, nil
}

func encodeInvokeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

func makeInvokeActionEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		params := httprouter.ParamsFromContext(ctx)
		name := params.ByName("name")

		payload := req.(invokeRequest)
		result, err := s.RunAction(name, payload)
		if err != nil {

		}
		return result, nil
	}
}
