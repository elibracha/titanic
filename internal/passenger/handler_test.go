package passenger

import (
	ctx "context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"titanic-api/pkg/histogram"
	"titanic-api/pkg/response"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

/*
Test objects
*/
var (
	mService *MockService
	handler  *Handler
)

// MockVulnerabilityService is a mocked object that implements an interface
// that describes an object that the code we're testing relies on.
type MockService struct {
	mock.Mock
}

func (ms *MockService) FarePercentileHistogram() (*histogram.Histogram, error) {
	args := ms.Called()
	var res *histogram.Histogram
	if args.Get(0) != nil {
		res = args.Get(0).(*histogram.Histogram)
	}
	return res, args.Error(1)
}

func (ms *MockService) Get(pid int) (*Passenger, error) {
	args := ms.Called(pid)
	var res *Passenger
	if args.Get(0) != nil {
		res = args.Get(0).(*Passenger)
	}
	return res, args.Error(1)
}

func (ms *MockService) GetAll() ([]*Passenger, error) {
	args := ms.Called()
	var res []*Passenger
	if args.Get(0) != nil {
		res = args.Get(0).([]*Passenger)
	}
	return res, args.Error(1)
}

// pre test setup function
func setup() {
	mService = new(MockService)
	handler = NewHandler(mService)
}

/*
Test functions
*/
func TestHandlerGetAll_ValidRequest_ResponseOk(t *testing.T) {
	setup()

	passengers := createPassengers(3)

	// given
	r, err := http.NewRequest("GET", "/passenger", nil)
	if err != nil {
		t.Fatal(err)
	}
	mService.On("GetAll").Return(createPassengers(3), nil /* error */)

	w := httptest.NewRecorder()

	// when
	handler.GetAll(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("Response As Expected", func() {
			var p []*Passenger
			err := json.NewDecoder(w.Body).Decode(&p)
			So(err, ShouldBeNil)
			So(len(p), ShouldEqual, len(passengers))
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGetAll_ValidRequest_ResponseInternalError(t *testing.T) {
	setup()

	// given
	r, err := http.NewRequest("GET", "/passenger", nil)
	if err != nil {
		t.Fatal(err)
	}
	mService.On("GetAll").Return(nil, errors.New("error"))

	w := httptest.NewRecorder()

	// when
	handler.GetAll(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Result As Expected", func() {
			var errorResp response.Error
			err := json.NewDecoder(w.Body).Decode(&errorResp)
			So(err, ShouldBeNil)
			So(errorResp.Message, ShouldEqual, response.ErrInternalFailure.Error())
			So(errorResp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_ValidRequest_ResponseOk(t *testing.T) {
	setup()

	passenger := createPassengers(1)[0]

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(passenger.PassengerId))
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	mService.On("Get", passenger.PassengerId).Return(passenger, nil /* error */)

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("Response As Expected", func() {
			var rs *Response
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.PassengerId, ShouldEqual, passenger.PassengerId)
			So(rs.Name, ShouldEqual, passenger.Name)
			So(rs.Age, ShouldEqual, passenger.Age)
			So(rs.Sex, ShouldEqual, passenger.Sex)
			So(rs.Cabin, ShouldEqual, passenger.Cabin)
			So(rs.Parch, ShouldEqual, passenger.Parch)
			So(rs.Pclass, ShouldEqual, passenger.Pclass)
			So(rs.Embarked, ShouldEqual, passenger.Embarked)
			So(rs.Survived, ShouldEqual, passenger.Survived)
			So(rs.Ticket, ShouldEqual, passenger.Ticket)
			So(rs.Fare, ShouldEqual, passenger.Fare)
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_ValidRequestWithAttributes_ResponseOk(t *testing.T) {
	setup()

	passenger := createPassengers(2)[1]

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(passenger.PassengerId))
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create a new url.Values and add query parameters
	queryParams := url.Values{}
	queryParams.Add("attributes", "id,name, age")

	// Set the query parameters as the RawQuery of the request URL
	r.URL.RawQuery = queryParams.Encode()

	mService.On("Get", passenger.PassengerId).Return(passenger, nil /* error */)

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("Response As Expected", func() {
			var rs *Response
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.PassengerId, ShouldEqual, passenger.PassengerId)
			So(rs.Name, ShouldEqual, passenger.Name)
			So(rs.Age, ShouldEqual, passenger.Age)
			So(rs.Sex, ShouldEqual, "")
			So(rs.Cabin, ShouldEqual, "")
			So(rs.Parch, ShouldEqual, 0)
			So(rs.Pclass, ShouldEqual, 0)
			So(rs.Embarked, ShouldEqual, "")
			So(rs.Survived, ShouldEqual, 0)
			So(rs.Ticket, ShouldEqual, "")
			So(rs.Fare, ShouldEqual, 0.0)
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_InvalidRequest_ResponseBadRequest(t *testing.T) {
	setup()

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid-id")
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("Response As Expected", func() {
			var rs *response.Error
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.Code, ShouldEqual, http.StatusBadRequest)
			So(rs.Message, ShouldEqual, ErrInvalidID.Error())
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_InvalidRequestWithAttributes_ResponseBadRequest(t *testing.T) {
	setup()

	passenger := createPassengers(2)[1]

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(passenger.PassengerId))
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create a new url.Values and add query parameters
	queryParams := url.Values{}
	queryParams.Add("attributes", "id,name, invalid-attribute")

	// Set the query parameters as the RawQuery of the request URL
	r.URL.RawQuery = queryParams.Encode()

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("Response As Expected", func() {
			var rs *response.Error
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.Code, ShouldEqual, http.StatusBadRequest)
			So(rs.Message, ShouldContainSubstring, "invalid-attribute")
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_InvalidRequestWithAttributesTwice_ResponseBadRequest(t *testing.T) {
	setup()

	passenger := createPassengers(2)[1]

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(passenger.PassengerId))
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create a new url.Values and add query parameters
	queryParams := url.Values{}
	queryParams.Add("attributes", "id,id")

	// Set the query parameters as the RawQuery of the request URL
	r.URL.RawQuery = queryParams.Encode()

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 400", func() {
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("Response As Expected", func() {
			var rs *response.Error
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.Code, ShouldEqual, http.StatusBadRequest)
			So(rs.Message, ShouldContainSubstring, "id")
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerGet_ValidRequest_ResponseInternalError(t *testing.T) {
	setup()

	passenger := createPassengers(2)[1]

	// given
	r, err := http.NewRequest("GET", "/passenger/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(passenger.PassengerId))
	r = r.WithContext(ctx.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	mService.On("Get", passenger.PassengerId).Return(nil, errors.New("error"))

	w := httptest.NewRecorder()

	// when
	handler.Get(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Response As Expected", func() {
			var rs *response.Error
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.Code, ShouldEqual, http.StatusInternalServerError)
			So(rs.Message, ShouldEqual, response.ErrInternalFailure.Error())
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerFareHistogram_ValidRequest_ResponseOk(t *testing.T) {
	setup()

	h := createHistogram()

	// given
	r, err := http.NewRequest("GET", "/passenger/fare/histogram/percentile", nil)
	if err != nil {
		t.Fatal(err)
	}
	mService.On("FarePercentileHistogram").Return(h, nil /* error */)

	w := httptest.NewRecorder()

	// when
	handler.FareHistogram(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("Response As Expected", func() {
			var hist *histogram.Histogram
			err := json.NewDecoder(w.Body).Decode(&hist)
			So(err, ShouldBeNil)
			So(len(hist.Entries), ShouldEqual, len(h.Entries))
		})
	})

	mService.AssertExpectations(t)
}

func TestHandlerFareHistogram_ValidRequest_ResponseInternalError(t *testing.T) {
	setup()

	// given
	r, err := http.NewRequest("GET", "/passenger/fare/histogram/percentile", nil)
	if err != nil {
		t.Fatal(err)
	}
	mService.On("FarePercentileHistogram").Return(nil, errors.New("error"))

	w := httptest.NewRecorder()

	// when
	handler.FareHistogram(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 500", func() {
			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Response As Expected", func() {
			var rs *response.Error
			err := json.NewDecoder(w.Body).Decode(&rs)
			So(err, ShouldBeNil)
			So(rs.Code, ShouldEqual, http.StatusInternalServerError)
			So(rs.Message, ShouldEqual, response.ErrInternalFailure.Error())
		})
	})

	mService.AssertExpectations(t)
}

func createPassengers(size int) []*Passenger {
	var passengers []*Passenger
	for i := 0; i < size; i++ {
		p := &Passenger{
			PassengerId: i,
			Survived:    1,
			Pclass:      1,
			Name:        "John Doe",
			Sex:         "Male",
			Age:         "30",
			SibSp:       1,
			Parch:       0,
			Ticket:      "A123",
			Fare:        50.25,
			Cabin:       "C123",
			Embarked:    "S",
		}
		passengers = append(passengers, p)
	}
	return passengers
}

func createHistogram() *histogram.Histogram {
	return &histogram.Histogram{
		Entries: []*histogram.Entry{
			{Bin: 25, Count: 100},
			{Bin: 50, Count: 101},
			{Bin: 75, Count: 102},
			{Bin: 100, Count: 103},
		}}
}
