package user_limit

type UpdateUserLimitRequest struct {
	Limit int `json:"limit" validate:"required,min=300"`
}
