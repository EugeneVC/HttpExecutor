package models

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var allowCreateMethods = map[string]bool{
	"GET": true,
	"POST":true,
}

type Task struct {
	ID          int64
	Method      string
	URL         string
	Header      http.Header
	RequestBody string

	StatusCode int
	Length int
}

func (t Task) ValidateCreateRequest() error {
	if _, ok := allowCreateMethods[t.Method]; !ok {
		return errors.New("Unknown method")
	}
	if _, err := url.Parse(t.URL); err != nil {
		return errors.New("Wrong URL")
	}

	return nil
}

func (t Task) String() string{
	return fmt.Sprintf("%d [%s] %s",t.ID,t.Method,t.URL)
}
