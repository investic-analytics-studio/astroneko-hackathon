package waiting_list

type JoinWaitingListRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CheckWaitingListRequest struct {
	Email string `json:"email" validate:"required,email"`
}
