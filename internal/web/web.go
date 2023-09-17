package web

import (
	"html/template"
	"log"
	"net/http"
	"titanic-api/internal/passenger"
	"titanic-api/pkg/histogram"
)

type Data struct {
    Passengers []*passenger.Passenger
    Histogram  []*histogram.Entry
}


type Handler struct {
    service passenger.Service
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("public/layout.html")
    if err != nil {
        log.Println(err.Error())
    }

    tmpl.Execute(w, nil)
}

func (h *Handler) Passengers(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("public/passengers.html")
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
    tmpl, err := template.ParseFiles("public/histogram.html")
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

