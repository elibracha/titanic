package response

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

var (
	ErrInternalFailure = errors.New("something went wrong, please try again")
)

// Error represents the structure of an error response.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SendError sends an error response in JSON format.
func SendError(r *http.Request, w http.ResponseWriter, code int, message string) {
	err := Error{
		Message: message,
		Code:    code,
	}
	w.WriteHeader(code)
	render.JSON(w, r, err)
}

// SendBody sends response in JSON format.
func SendBody(r *http.Request, w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	render.JSON(w, r, body)
}
