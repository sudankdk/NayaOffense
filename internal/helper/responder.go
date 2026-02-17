package helper

import (
	"context"
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
