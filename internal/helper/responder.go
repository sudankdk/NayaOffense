package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sudankdk/offense/internal/domain"
	"github.com/sudankdk/offense/internal/tools"
)

func ExecuteResponse(ctx context.Context, tools tools.Tool, input map[string]string) domain.Response {
	start := time.Now()
	resp := domain.Response{
		Tool: tools.Name(),
		Meta: domain.Meta{
			StartedAt: start,
		},
	}
	data, err := tools.Run(ctx, input)
	resp.Meta.Duration = time.Since(start).Milliseconds()
	if err != nil {
		resp.Status = "error"
		resp.Error = &domain.ErrorPayload{
			Message: err.Error(),
			Code:    "TOOL_EXECUTION_FAILED",
		}
		return resp
	}
	resp.Status = "success"
	resp.Data = data
	return resp
}

// ExecuteResponseStream executes a tool and streams events in SSE format
func ExecuteResponseStream(ctx context.Context, tool tools.Tool, input map[string]string, w io.Writer, flusher http.Flusher) {
	start := time.Now()

	// Send start event
	startEvent := map[string]interface{}{
		"type":      "start",
		"tool":      tool.Name(),
		"timestamp": start,
		"input":     input,
	}
	writeSSE(w, startEvent)
	flusher.Flush()

	// Validate input
	if validator, ok := tool.(interface{ Validate(map[string]string) error }); ok {
		if err := validator.Validate(input); err != nil {
			errorEvent := map[string]interface{}{
				"type":      "error",
				"timestamp": time.Now(),
				"error": map[string]string{
					"message": err.Error(),
					"code":    "VALIDATION_FAILED",
				},
			}
			writeSSE(w, errorEvent)
			flusher.Flush()
			return
		}
	}

	// Execute tool
	data, err := tool.Run(ctx, input)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		// Send error event
		errorEvent := map[string]interface{}{
			"type":      "error",
			"timestamp": time.Now(),
			"duration":  duration,
			"error": map[string]string{
				"message": err.Error(),
				"code":    "TOOL_EXECUTION_FAILED",
			},
		}
		writeSSE(w, errorEvent)
		flusher.Flush()
		return
	}

	// Send success event with data
	successEvent := map[string]interface{}{
		"type":      "complete",
		"tool":      tool.Name(),
		"timestamp": time.Now(),
		"duration":  duration,
		"status":    "success",
		"data":      data,
	}
	writeSSE(w, successEvent)
	flusher.Flush()
}

// writeSSE writes data in Server-Sent Events format
func writeSSE(w io.Writer, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Fallback error message
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"failed to encode event\"}\n\n")
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
}
