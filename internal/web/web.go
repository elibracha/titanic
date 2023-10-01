package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"titanic-api/internal/passenger"
	"titanic-api/pkg/histogram"

	"github.com/go-chi/chi"
)

type Data struct {
	Passengers []*passenger.Passenger
	Histogram  []*histogram.Entry
}

type Handler struct {
	service passenger.Service
}

func (h *Handler) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", h.Root)
	router.Get("/passengers", h.Passengers)
	router.Get("/passenger/{id}", h.Passenger)
	router.Get("/histogram", h.Histogram)
	return router
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		log.Println(err.Error())
	}

	tmpl.Execute(w, nil)
}

func (h *Handler) Passenger(w http.ResponseWriter, r *http.Request) {
	var data Data
	pid, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}

	p, err := h.service.Get(pid)
	if err == nil {
		data.Passengers = []*passenger.Passenger{p}
	}

	var tmpl *template.Template
	switch len(data.Passengers) {
	case 0:
		tmpl, err = template.ParseFiles("templates/404.html")
	default:
		tmpl, err = template.ParseFiles("templates/passengers.html")
	}

	if err != nil {
		log.Println(err.Error())
	}

	tmpl.Execute(w, data)
}

func (h *Handler) Passengers(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/passengers.html")
	if err != nil {
		log.Println(err.Error())
	}

	var data Data
	p, err := h.service.GetAll()
	if err == nil {
		data.Passengers = p
	}

	tmpl.Execute(w, data)
}

func (h *Handler) Histogram(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/histogram.html")
	if err != nil {
		log.Println(err.Error())
	}

	var data Data
	pc, err := h.service.FarePercentileHistogram()
	if err == nil {
		data.Histogram = pc.Entries
	}

	tmpl.Execute(w, data)
}

func NewHandler(service passenger.Service) *Handler {
	return &Handler{service: service}
}
