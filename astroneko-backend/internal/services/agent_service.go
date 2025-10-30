package services

import (
	"context"

	"astroneko-backend/internal/core/domain/agent"
	agentPorts "astroneko-backend/internal/core/ports/agent"
	"astroneko-backend/pkg/logger"
)

type AgentService struct {
	agentRepo agentPorts.RepositoryInterface
	logger    logger.Logger
}

func NewAgentService(agentRepo agentPorts.RepositoryInterface, log logger.Logger) *AgentService {
	return &AgentService{
		agentRepo: agentRepo,
		logger:    log,
	}
}

func (s *AgentService) ClearState(ctx context.Context, userID string, request agent.ClearStateRequest) (*agent.ClearStateResponse, error) {

	s.logger.Info("Clearing agent state for user",
		logger.Field{Key: "module", Value: "agent_service"},
		logger.Field{Key: "user_id", Value: userID})

	response, err := s.agentRepo.ClearState(ctx, request)
	if err != nil {
		s.logger.Error("Failed to clear agent state",
			logger.Field{Key: "module", Value: "agent_service"},
			logger.Field{Key: "user_id", Value: userID},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	s.logger.Info("Agent state cleared successfully",
		logger.Field{Key: "module", Value: "agent_service"},
		logger.Field{Key: "user_id", Value: userID},
		logger.Field{Key: "status", Value: response.Status})

	return response, nil
}

func (s *AgentService) Reply(ctx context.Context, userID string, request agent.ReplyRequest) (*agent.ReplyResponse, error) {
	s.logger.Info("Sending message to agent",
		logger.Field{Key: "module", Value: "agent_service"},
		logger.Field{Key: "user_id", Value: userID},
		logger.Field{Key: "text_length", Value: len(request.Text)})

	response, err := s.agentRepo.Reply(ctx, request)
	if err != nil {
		s.logger.Error("Failed to get agent reply",
			logger.Field{Key: "module", Value: "agent_service"},
			logger.Field{Key: "user_id", Value: userID},
			logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	s.logger.Info("Agent reply received successfully",
		logger.Field{Key: "module", Value: "agent_service"},
		logger.Field{Key: "user_id", Value: userID},
		logger.Field{Key: "status", Value: response.Status},
		logger.Field{Key: "message_length", Value: len(response.Message)})

	return response, nil
}
