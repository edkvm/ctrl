package http

import (
	"context"
	"github.com/edkvm/ctrl/pkg/endpoint"
	"net/http"
)

type EncodeResponseFunc func(context.Context, http.ResponseWriter, interface{}) error
type DecodeRequestFunc func(context.Context, *http.Request) (request interface{}, err error)

type server struct {
	e endpoint.Endpoint
	dec DecodeRequestFunc
	enc EncodeResponseFunc
}

func NewServer(e endpoint.Endpoint, dec DecodeRequestFunc, enc EncodeResponseFunc) *server {
	return &server{
		e: e,
		dec: dec,
		enc: enc,
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec, err := s.dec(ctx, r)
	if err != nil {
		return
	}

	result, err := s.e(ctx, dec)
	if err != nil {
		return
	}

	if err := s.enc(ctx, w, result); err != nil {
		return
	}


}
