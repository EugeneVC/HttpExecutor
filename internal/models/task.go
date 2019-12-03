package models

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var allowMethods = map[string]bool{
	"GET": true,
	"POST":true,
	"LIST" :true,
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

func (t Task) ValidateRequest() error {
	if _, ok := allowMethods[t.Method]; !ok {
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
