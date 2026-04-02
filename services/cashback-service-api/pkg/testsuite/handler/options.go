package handler

import (
	"encoding/json"
	"fmt"
	"strings"
)

func WithJSONBody(body string) OptionFunc {
	return func(d *requestData) {
		d.body = strings.NewReader(body)
		d.headers["Content-Type"] = "application/json"
	}
}

func WithJSONBodyStruct(v any) OptionFunc {
	b, _ := json.Marshal(v)
	return WithJSONBody(string(b))
}

func WithPathParam(name string, value any) OptionFunc {
	return func(d *requestData) {
		if d.params == nil {
			d.params = map[string]string{}
		}
		d.params[name] = fmt.Sprintf("%v", value)
	}
}

func WithQuery(q string) OptionFunc {
	return func(d *requestData) {
		d.query = q
	}
}
