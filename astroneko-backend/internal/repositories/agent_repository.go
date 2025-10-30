package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"astroneko-backend/configs"
	"astroneko-backend/internal/core/domain/agent"
	agentPorts "astroneko-backend/internal/core/ports/agent"
	"astroneko-backend/pkg/apprequest"
)

type agentRepository struct {
	agentBaseURL string
	token        string
	httpClient   apprequest.HTTPRequest
}

func NewAgentRepository() agentPorts.RepositoryInterface {
	return &agentRepository{
		agentBaseURL: configs.GetViper().ExternalURL.AstronekoURL,
		token:        configs.GetViper().ExternalURL.Token,
		httpClient:   apprequest.NewRequester(),
	}
}

func (r *agentRepository) ClearState(ctx context.Context, request agent.ClearStateRequest) (*agent.ClearStateResponse, error) {
	uri := fmt.Sprintf("%s/api/cat-fortune/clear-state", r.agentBaseURL)

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, resp := r.httpClient.NewRequest(payload, apprequest.DELETE, uri)
	req.Header.SetContentTypeBytes(apprequest.ApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+r.token)

	{
		err = fasthttp.Do(req, resp)
		if err != nil {
			logrus.Error("error on clear-state request: ", err)
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
	}

	{
		body := resp.Body()
		var response agent.ClearStateResponse

		if resp.StatusCode() != fasthttp.StatusOK {
			logrus.Printf("Clear-state API error - status: %d, body: %s", resp.StatusCode(), string(body))

			// Handle authentication errors specifically
			if resp.StatusCode() == fasthttp.StatusUnauthorized {
				return nil, fmt.Errorf("authentication failed: invalid or expired token")
			}

			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			logrus.Error("error on unmarshal clear-state response: ", err)
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return &response, nil
	}
}

func (r *agentRepository) Reply(ctx context.Context, request agent.ReplyRequest) (*agent.ReplyResponse, error) {
	uri := fmt.Sprintf("%s/api/cat-fortune/reply", r.agentBaseURL)

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, resp := r.httpClient.NewRequest(payload, apprequest.POST, uri)
	req.Header.SetContentTypeBytes(apprequest.ApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+r.token)

	{
		err = fasthttp.Do(req, resp)
		if err != nil {
			logrus.Error("error on reply request: ", err)
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
	}

	{
		body := resp.Body()
		var apiResponse agent.ReplyResponseFromAPI

		if resp.StatusCode() != fasthttp.StatusOK {
			logrus.Printf("Reply API error - status: %d, body: %s", resp.StatusCode(), string(body))

			// Handle authentication errors specifically
			if resp.StatusCode() == fasthttp.StatusUnauthorized {
				return nil, fmt.Errorf("authentication failed: invalid or expired token")
			}

			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
		}

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			logrus.Error("error on unmarshal reply response: ", err)
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return apiResponse.ToReplyResponse(), nil
	}
}
