package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type (
	Response struct {
		Code int
		Body string
	}

	requestData struct {
		body    io.Reader
		query   string
		headers map[string]string
		params  map[string]string
	}

	OptionFunc func(*requestData)

	Suite struct {
		suite.Suite
		router   chi.Router
		recorder *httptest.ResponseRecorder
		method   string
	}
)

func (s *Suite) PrepareRouter(method, pattern string, h http.HandlerFunc) {
	s.router = chi.NewRouter()
	s.recorder = httptest.NewRecorder()
	s.method = method
	s.router.Method(method, pattern, h)
}

func (s *Suite) Serve(path string, opts ...OptionFunc) {
	d := &requestData{headers: map[string]string{"Content-Type": "application/json"}}
	for _, opt := range opts {
		opt(d)
	}

	target := path
	for k, v := range d.params {
		target = strings.ReplaceAll(target, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	if d.query != "" {
		target += "?" + d.query
	}

	req := httptest.NewRequest(s.method, target, d.body)
	for k, v := range d.headers {
		req.Header.Set(k, v)
	}

	s.router.ServeHTTP(s.recorder, req)
}

func (s *Suite) Response() Response {
	return Response{
		Code: s.recorder.Code,
		Body: s.recorder.Body.String(),
	}
}
