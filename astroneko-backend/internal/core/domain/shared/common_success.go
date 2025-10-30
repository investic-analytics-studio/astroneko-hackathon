package shared

import "net/http"

type CustomSuccess struct {
	HTTPStatus int
	Code       string
	Message    string
}

var successRegistry = map[string]Status{
	"SUC_200": {
		HTTPStatus: http.StatusOK,
		Code:       "SUC_200",
		Message:    []string{"OK", "Standard successful GET"},
	},
	"SUC_201": {
		HTTPStatus: http.StatusCreated,
		Code:       "SUC_201",
		Message:    []string{"Created", "Resource created successfully"},
	},
	"SUC_202": {
		HTTPStatus: http.StatusAccepted,
		Code:       "SUC_202",
		Message:    []string{"Accepted", "Request accepted for processing"},
	},
	"SUC_204": {
		HTTPStatus: http.StatusNoContent,
		Code:       "SUC_204",
		SuccessID:  "SUC_204",
		Message:    []string{"No Content", "Successful with no return data"},
	},
}

func NewSuccessResponse(code string, detailOverride ...string) (int, ResponseBody) {
	status, ok := successRegistry[code]
	if !ok {
		status = successRegistry["SUC_200"]
	}

	if len(detailOverride) > 0 {
		status.Message = detailOverride
	}

	return status.HTTPStatus, ResponseBody{
		Status: status,
	}
}
