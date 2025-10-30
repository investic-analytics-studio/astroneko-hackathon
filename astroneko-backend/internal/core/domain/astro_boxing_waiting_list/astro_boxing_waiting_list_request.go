package astro_boxing_waiting_list

type JoinAstroBoxingWaitingListRequest struct {
	Email string `json:"email" validate:"required,email"`
}
