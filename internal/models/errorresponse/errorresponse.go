package errorresponse

import (
	"net/http"
)

var exceptionStatusNameMapping = map[uint]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusForbidden:           "Forbidden",
}

// Exception struct
type Exception struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  uint   `json:"status"`
}

// ExceptionResponse response similar to existing exception
func ExceptionResponse(status uint, message string) Exception {
	return Exception{
		Name:    exceptionStatusNameMapping[status],
		Message: message,
		Code:    0,
		Status:  status,
	}
}
