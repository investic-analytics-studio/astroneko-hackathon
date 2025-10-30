package agent

import (
	"context"

	"astroneko-backend/internal/core/domain/agent"
)

type ServiceInterface interface {
	ClearState(ctx context.Context, userID string, request agent.ClearStateRequest) (*agent.ClearStateResponse, error)
	Reply(ctx context.Context, userID string, request agent.ReplyRequest) (*agent.ReplyResponse, error)
}
