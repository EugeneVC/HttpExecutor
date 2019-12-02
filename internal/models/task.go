package models

import (
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

func (t Task) ValidateRequest() bool {
	if _, ok := allowMethods[t.Method]; !ok {
		return false
	}
	if _, err := url.Parse(t.URL); err != nil {
		return false
	}

	return true
}

func (t Task) Generate() {

}
