package agent

type ClearStateResponse struct {
	Status string `json:"status"`
}

type ReplyResponseFromAPI struct {
	Status    string `json:"status"`
	Text      string `json:"text"`
	Card      string `json:"card,omitempty"`
	Meaning   string `json:"meaning,omitempty"`
	SessionID string `json:"session_id"`
}

type ReplyResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Card      string `json:"card,omitempty"`
	Meaning   string `json:"meaning,omitempty"`
	SessionID string `json:"session_id"`
}

func (r *ReplyResponseFromAPI) ToReplyResponse() *ReplyResponse {
	return &ReplyResponse{
		Status:    r.Status,
		Message:   r.Text,
		Card:      r.Card,
		Meaning:   r.Meaning,
		SessionID: r.SessionID,
	}
}
