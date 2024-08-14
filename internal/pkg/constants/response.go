package constants

const (
	STATUS_SUCCESS        = "200"
	STATUS_BAD_REQUEST    = "400"
	STATUS_UNAUTHORIZED   = "401"
	STATUS_FORBIDDEN      = "403"
	STATUS_DATA_NOT_FOUND = "404"
	STATUS_CONFLICT       = "409"
	STATUS_GENERAL_ERROR  = "500"
)

const (
	MESSAGE_SUCCESS        = "success"
	MESSAGE_SUCCESS_CREATE = "success create data"
	MESSAGE_SUCCESS_UPDATE = "success update data"
	MESSAGE_SUCCESS_DELETE = "success delete data"
	MESSAGE_BAD_REQUEST    = "invalid request format"
	MESSAGE_DATA_NOT_FOUND = "data not found"
	MESSAGE_FAILED         = "something went wrong"
)

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationData struct {
	Page        uint `json:"page"`
	TotalPages  uint `json:"totalPages"`
	TotalItems  uint `json:"totalItems"`
	Limit       uint `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

type PaginationResponseData struct {
	Results        interface{} `json:"results"`
	PaginationData `json:"pagination"`
}

func ErrorResponse(status, message string) (response DefaultResponse) {
	if message == "" {
		message = MESSAGE_FAILED
	}
	response = DefaultResponse{
		Status:  status,
		Message: message,
		Data:    struct{}{},
	}
	return
}
