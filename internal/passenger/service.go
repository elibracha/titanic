package passenger

import (
	"titanic-api/pkg/histogram"
)

type Service interface {
	Get(pid int) (*Passenger, error)
	GetAll() ([]*Passenger, error)
	FarePercentileHistogram() (*histogram.Histogram, error)
}

type service struct {
	store Store
}

func (s *service) FarePercentileHistogram() (*histogram.Histogram, error) {
	passengers, err := s.store.GetPassengers()
	if err != nil {
		return nil, err
	}

	var fares []float64
	for _, p := range passengers {
		fares = append(fares, p.Fare)
	}

	return histogram.Percentile(fares), nil
}

func (s *service) Get(pid int) (*Passenger, error) {
	return s.store.GetPassenger(pid)
}

func (s *service) GetAll() ([]*Passenger, error) {
	return s.store.GetPassengers()
}

func NewService(store Store) Service {
	return &service{store: store}
}
