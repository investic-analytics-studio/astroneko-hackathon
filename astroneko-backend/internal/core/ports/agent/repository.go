package agent

import (
	"context"

	"astroneko-backend/internal/core/domain/agent"
)

type RepositoryInterface interface {
	ClearState(ctx context.Context, request agent.ClearStateRequest) (*agent.ClearStateResponse, error)
	Reply(ctx context.Context, request agent.ReplyRequest) (*agent.ReplyResponse, error)
}
