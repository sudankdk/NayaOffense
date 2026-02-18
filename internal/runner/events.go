package runner

import "time"

type EventType string

const (
	EventStart    EventType = "start"
	EventStdout   EventType = "stdout"
	EventStderr   EventType = "stderr"
	EventError    EventType = "error"
	EventComplete EventType = "complete"
)

type Event struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      string    `json:"data,omitempty"`
}
