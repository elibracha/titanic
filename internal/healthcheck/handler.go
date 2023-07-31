package healthcheck

import (
	"github.com/go-chi/chi"
	"net/http"
	"titanic-api/pkg/response"
)

type Handler struct{}

type Status struct {
	Code  int    `json:"code"`
	State string `json:"state"`
}

func (h *Handler) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", h.Health)
	return router
}

// Package 	godoc
// @Summary Get health check status
// @Description Get health check status
// @Tags    health
// @ID 		healthcheck
// @Produce json
// @Success 200 {object} Status
// @Router  /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	status := &Status{
		Code:  http.StatusOK,
		State: "OK",
	}
	response.SendBody(r, w, http.StatusOK, status)
}

func NewHandler() *Handler {
	return &Handler{}
}
