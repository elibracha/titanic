package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
Test objects
*/
var (
	handler *Handler
)

// pre test setup function
func setup() {
	handler = NewHandler()
}

func TestHandler_ValidRequest_ResponseOk(t *testing.T) {
	setup()

	status := &Status{
		Code:  200,
		State: "OK",
	}

	// given
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// when
	handler.Health(w, r)

	// then
	Convey("Test handler\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("Response As Expected", func() {
			var s *Status
			err := json.NewDecoder(w.Body).Decode(&s)
			So(err, ShouldBeNil)
			So(s.Code, ShouldEqual, status.Code)
			So(s.State, ShouldEqual, status.State)
		})
	})
}
