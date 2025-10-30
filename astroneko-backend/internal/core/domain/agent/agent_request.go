package agent

type ClearStateRequest struct {
	SessionID string `json:"session_id"`
}

type ReplyRequest struct {
	Text      string `json:"text" validate:"required"`
	UserID    string `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}
