package waiting_list

import (
	"astroneko-backend/internal/core/domain/shared"
)

type WaitingListUserResponse struct {
	shared.EmailResponse
}

func (w *WaitingListUser) ToResponse() *WaitingListUserResponse {
	return &WaitingListUserResponse{
		EmailResponse: *shared.NewEmailResponse(
			w.ID.String(),
			w.Email,
			w.CreatedAt,
			w.UpdatedAt,
		),
	}
}

type JoinWaitingListResponse struct {
	shared.ResponseBody
}

type IsInWaitingListResponse struct {
	shared.StatusResponse
}

func NewIsInWaitingListResponse(isInWaitingList bool) *IsInWaitingListResponse {
	message := "User is not in waiting list"
	if isInWaitingList {
		message = "User is in waiting list"
	}
	return &IsInWaitingListResponse{
		StatusResponse: *shared.NewStatusResponse(true, message),
	}
}
