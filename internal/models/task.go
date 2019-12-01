package models

import "net/http"

type TaskMethod uint8

type Task struct {
	Method TaskMethod
	URL    string
	Header http.Handler
	Body   string
}
