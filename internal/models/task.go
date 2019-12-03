package models

import (
	"errors"
	"net/http"
	"net/url"
)

var allowMethods = map[string]bool{
	"GET": true,
}

type Task struct {
	ID          int64
	Method      string
	URL         string
	Header      http.Handler
	RequestBody string
}

func (t Task) ValidateRequest() error {
	if _, ok := allowMethods[t.Method]; !ok {
		return errors.New("Unknown method")
	}
	if _, err := url.Parse(t.URL); err != nil {
		return errors.New("Wrong URL")
	}

	return nil
}

func (t Task) Generate() {

}
