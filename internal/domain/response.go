package domain

import "time"

type Response struct {
	Tool   string        `json:"tool"`
	Data   any           `json:"data,omitempty"`
	Status string        `json:"status"`
	Meta   Meta          `json:"meta"`
	Error  *ErrorPayload `json:"error,omitempty"`
}

type Meta struct {
	StartedAt time.Time `json:"started_at"`
	Duration  int64     `json:"duration_ms"`
}

type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
