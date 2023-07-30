package passenger

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
)

type csvStore struct {
	path string
}

func (s *csvStore) GetPassenger(pid int) (*Passenger, error) {
	passengers, err := s.loadPassengers()
	if err != nil {
		return nil, err
	}

	for _, p := range passengers {
		if p.PassengerId == pid {
			return p, nil
		}
	}

	return nil, ErrPassengerNotFound
}

func (s *csvStore) GetPassengers() ([]*Passenger, error) {
	return s.loadPassengers()
}

func (s *csvStore) loadPassengers() ([]*Passenger, error) {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error opening store path: %s error: %s", s.path, err.Error())
	}
	defer file.Close()

	var passengers []*Passenger
	if err = gocsv.UnmarshalFile(file, &passengers); err != nil {
		return nil, fmt.Errorf("error loading store data path: %s error: %s", s.path, err.Error())
	}

	return passengers, nil
}

func NewStoreCSV(path string) Store {
	return &csvStore{path: path}
}
