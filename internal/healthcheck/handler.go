package healthcheck

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type HeathHandler struct{}

// Package 	godoc
// @Summary Get health check status
// @Description Get health check status
// @Tags    health
// @ID 		healthcheck
// @Produce json
// @Success 200
// @Router  /health [get]
func (h *HeathHandler) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		render.NoContent(w, r)
	})
	return router
}

func NewHandler() *HeathHandler {
	return &HeathHandler{}
}
