package passenger

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"titanic-api/pkgs/response"
)

const (
	maxAttributeParamLength = 256
)

var (
	ErrInvalidID = fmt.Errorf("id provided is not a valid integer")
)

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
	router.Get("/fare/histogram/percentile", h.FarePercentiles)
	return router
}

// Package 	godoc
// @Summary Get passengers
// @Description Get all passengers
// @Tags    passenger
// @ID 		passenger-get-all
// @Produce json
// @Success 200 {object} []Response
// @Failure 500 {object} response.Error
// @Router  /passenger [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	passengers, err := h.service.GetAll()
	switch err {
	case nil:
		var rs []*Response
		for _, p := range passengers {
			rs = append(rs, h.convertPassenger(p))
		}
		response.SendBody(r, w, http.StatusOK, rs)
	default:
		log.Println(fmt.Sprintf("request id: %s failed to get passengers: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		response.SendError(r, w, http.StatusInternalServerError, response.ErrInternalFailure.Error())
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
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router  /passenger/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	attr := r.URL.Query().Get("attributes")
	if err := h.validateAttributes(attr); err != nil {
		response.SendError(r, w, http.StatusBadRequest, err.Error())
		return
	}

	pid, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.SendError(r, w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	storePassenger, err := h.service.Get(pid)
	switch {
	case err == ErrPassengerNotFound:
		response.SendError(r, w, http.StatusNotFound, ErrPassengerNotFound.Error())
		return
	case err != nil:
		log.Println(fmt.Sprintf("request id: %s failed to get passenger: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		response.SendError(r, w, http.StatusInternalServerError, response.ErrInternalFailure.Error())
		return
	}

	p := h.convertPassenger(storePassenger)
	switch len(attr) {
	case 0:
		response.SendBody(r, w, http.StatusOK, p)
	default:
		response.SendBody(r, w, http.StatusOK, h.filterAttributes(p, attr))
	}
}

// Package 	godoc
// @Summary Get fare histogram histogram
// @Description Get histogram represention of number of passengers in each precentile
// @Tags    passenger
// @ID 		passenger-fare-histogram
// @Produce json
// @Success 200 {object} histogram.Histogram
// @Failure 500 {object} response.Error
// @Router  /passenger/fare/histogram/percentile [get]
func (h *Handler) FarePercentiles(w http.ResponseWriter, r *http.Request) {
	histogram, err := h.service.FarePercentileHistogram()
	switch err {
	case nil:
		response.SendBody(r, w, http.StatusOK, histogram)
	default:
		log.Println(fmt.Sprintf("request id: %s failed to get fare histogram histogram: %v",
			middleware.GetReqID(r.Context()), err.Error()))
		response.SendError(r, w, http.StatusInternalServerError, response.ErrInternalFailure.Error())
	}
}

func (h *Handler) validateAttributes(attr string) error {
	if len(attr) == 0 {
		return nil
	}

	attributes := strings.Split(attr, ",")
	if len(attributes) > maxAttributeParamLength {
		return fmt.Errorf("attributes query parameter is too long max length %d", maxAttributeParamLength)
	}

	for i, a := range attributes {
		attributes[i] = strings.TrimSpace(a)
	}

	passengerValue := reflect.ValueOf(&Response{}).Elem()
	passengerType := passengerValue.Type()

	visited := make(map[string]bool)
	for _, attribute := range attributes {
		var valid bool
		for i := 0; i < passengerValue.NumField(); i++ {
			if valid {
				break
			}
			fieldType := passengerType.Field(i)

			jsonField := fieldType.Tag.Get("json")
			switch jsonField {
			case attribute:
				if _, found := visited[attribute]; found {
					return fmt.Errorf("attribute '%s' provided multiple times in query", attribute)
				}
				visited[attribute] = true
				valid = true
			default:
				valid = false
			}
		}
		if !valid {
			return fmt.Errorf("unknown attribute filter provided '%s' in query", attribute)
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
