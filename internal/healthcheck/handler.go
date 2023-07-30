package healthcheck

import (
	"github.com/go-chi/chi"
	"net/http"
	"titanic-api/pkg/response"
)

type HealthHandler struct{}

type Status struct {
	Code  int    `json:"code"`
	State string `json:"state"`
}

// Package 	godoc
// @Summary Get health check status
// @Description Get health check status
// @Tags    health
// @ID 		healthcheck
// @Produce json
// @Success 200 {object} Status
// @Router  /health [get]
func (h *HealthHandler) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		status := &Status{
			Code:  http.StatusOK,
			State: "OK",
		}
		response.SendBody(r, w, http.StatusOK, status)
	})
	return router
}

func NewHandler() *HealthHandler {
	return &HealthHandler{}
}
