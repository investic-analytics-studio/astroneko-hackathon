package shared

// ResponseBody represents a response body.
type ResponseBody struct {
	Status                   Status                 `json:"status,omitempty"`
	Data                     interface{}            `json:"data,omitempty"`
	Meta                     map[string]interface{} `json:"meta,omitempty"`
	LatestAiSuggestQuestions *[]string              `json:"latestAiSuggestQuestions,omitempty"`

	CurrentPage *int   `json:"currentPage,omitempty"`
	PerPage     *int   `json:"perPage,omitempty"`
	TotalItem   *int64 `json:"totalItem,omitempty"`
}

// Status represents a status.
type Status struct {
	HTTPStatus int      `json:"-"`
	Code       string   `json:"code,omitempty"`
	SuccessID  string   `json:"successID,omitempty"`
	ErrorID    string   `json:"errorID,omitempty"`
	Message    []string `json:"message,omitempty"`
}
