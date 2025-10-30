package shared

// Pagination represents pagination parameters for API requests.
type Pagination struct {
	Limit  int `json:"limit" query:"limit" validate:"gte=-1,lte=100"`
	Offset int `json:"offset" query:"offset"`
}

// SortMethod represents sorting parameters for API requests.
type SortMethod struct {
	Asc     bool   `json:"asc" query:"asc"`
	OrderBy string `json:"orderBy" query:"orderBy"`
}
