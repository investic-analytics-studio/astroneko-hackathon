package user_limit

import "time"

type UserLimitResponse struct {
	ID        string    `json:"id"`
	Limit     int       `json:"limit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserLimitResponse struct {
	UserLimitResponse
}

type UpdateUserLimitResponse struct {
	UserLimitResponse
}

type IsUserOverLimitUsedResponse struct {
	IsOverLimit bool `json:"is_over_limit"`
}
