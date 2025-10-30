package astro_boxing_waiting_list

import (
	"time"

	"astroneko-backend/internal/core/domain/shared"
)

type AstroBoxingWaitingListUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AstroBoxingWaitingListUser) ToResponse() *AstroBoxingWaitingListUserResponse {
	return &AstroBoxingWaitingListUserResponse{
		ID:        a.ID.String(),
		Email:     a.Email,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type JoinAstroBoxingWaitingListResponse struct {
	shared.ResponseBody
}

type IsInAstroBoxingWaitingListResponse struct {
	IsInWaitingList bool `json:"is_in_waiting_list"`
}
