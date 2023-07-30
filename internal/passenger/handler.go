package passenger

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const maxAttributeParamLength = 256

type Response struct {
	PassengerId int     `json:"id"`
	Survived    int     `json:"survived"`
	Pclass      int     `json:"class"`
	Name        string  `json:"name"`
	Sex         string  `json:"sex"`
	Age         string  `json:"age"`
	SibSp       int     `json:"siblings-spouses"`
	Parch       int     `json:"parents-children"`
	Ticket      string  `json:"ticket"`
	Fare        float64 `json:"fare"`
	Cabin       string  `json:"cabin"`
	Embarked    string  `json:"embarked"`
}

type Handler struct {
	service Service
}

func (h *Handler) RegisterHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", h.GetAll)
	router.Get("/{id}", h.Get)
	router.Get("/fare/histogram/histogram", h.FarePercentiles)
	return router
}

// Package 	godoc
// @Summary Get fare histogram histogram
// @Description Get histogram represention of number of passengers in each precentile
// @Tags    passenger
// @ID 		passenger-fare-histogram
// @Produce json
// @Success 200 {object} histogram.Histogram
// @Failure 500
// @Router  /passenger/fare/histogram/histogram [get]
func (h *Handler) FarePercentiles(w http.ResponseWriter, r *http.Request) {
	histogram, err := h.service.FarePercentileHistogram()
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, histogram)
	default:
		log.Println(fmt.Sprintf("request id: %s failed to get fare histogram histogram: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.NoContent(w, r)
	}
}

// Package 	godoc
// @Summary Get passenger
// @Description Get passenger by ID number
// @Tags    passenger
// @ID 		passenger-get
// @Produce json
// @Param id path int true "Passenger ID"
// @Param attributes query []string false "Allowed: id, age, sex, name, survived, class, siblings-spouses, parents-children, ticket, fare, cabin, embarked"
// @Success 200 {object} Response
// @Failure 500
// @Router  /passenger/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	attr := r.URL.Query().Get("attributes")
	if err := h.validateAttributes(attr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.NoContent(w, r)
		return
	}

	pid, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.NoContent(w, r)
		return
	}

	storePassenger, err := h.service.Get(pid)
	if err != nil {
		log.Println(fmt.Sprintf("request id: %s failed to get passenger: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.NoContent(w, r)
		return
	}

	p := h.convertPassenger(storePassenger)
	switch len(attr) {
	case 0:
		render.JSON(w, r, p)
	default:
		render.JSON(w, r, h.filterAttributes(p, attr))
	}

	w.WriteHeader(http.StatusOK)
}

// Package 	godoc
// @Summary Get passengers
// @Description Get all passengers
// @Tags    passenger
// @ID 		passenger-get-all
// @Produce json
// @Success 200 {object} []Response
// @Failure 500
// @Router  /passenger [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	passengers, err := h.service.GetAll()
	switch err {
	case nil:
		var rs []*Response
		for _, p := range passengers {
			rs = append(rs, h.convertPassenger(p))
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, rs)
	default:
		log.Println(fmt.Sprintf("request id: %s failed to get passengers: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.NoContent(w, r)
	}
}

func (h *Handler) validateAttributes(attr string) error {
	attributes := strings.Split(attr, ",")
	if len(attributes) == 0 {
		return nil
	}
	if len(attributes) > maxAttributeParamLength {
		return fmt.Errorf("attributes query parameter is too long max length %d", maxAttributeParamLength)
	}

	for i, a := range attributes {
		attributes[i] = strings.TrimSpace(a)
	}

	passengerValue := reflect.ValueOf(&Response{}).Elem()
	passengerType := passengerValue.Type()

	for _, attr := range attributes {
		var valid bool
		for i := 0; i < passengerValue.NumField(); i++ {
			if valid {
				break
			}
			fieldType := passengerType.Field(i)

			jsonField := fieldType.Tag.Get("json")
			switch jsonField {
			case attr:
				valid = true
			default:
				valid = false
			}
		}
		if !valid {
			return fmt.Errorf("unknown attribute filter provided %s", attr)
		}
	}
	return nil
}

func (h *Handler) filterAttributes(p *Response, attr string) map[string]interface{} {
	attributes := strings.Split(attr, ",")
	rs := make(map[string]interface{})

	for i, a := range attributes {
		attributes[i] = strings.TrimSpace(a)
	}

	passengerValue := reflect.ValueOf(p).Elem()
	passengerType := passengerValue.Type()

	for _, a := range attributes {
		for i := 0; i < passengerValue.NumField(); i++ {
			fieldValue := passengerValue.Field(i)
			fieldType := passengerType.Field(i)

			jsonField := fieldType.Tag.Get("json")
			if jsonField != a {
				continue
			}
			rs[jsonField] = fieldValue.Interface()
		}
	}
	return rs
}

func (h *Handler) convertPassenger(source *Passenger) *Response {
	var dest Response
	// Directly assign values from source to destination
	dest.PassengerId = source.PassengerId
	dest.Survived = source.Survived
	dest.Pclass = source.Pclass
	dest.Name = source.Name
	dest.Sex = source.Sex
	dest.Age = source.Age
	dest.SibSp = source.SibSp
	dest.Parch = source.Parch
	dest.Ticket = source.Ticket
	dest.Fare = source.Fare
	dest.Cabin = source.Cabin
	dest.Embarked = source.Embarked

	return &dest
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
