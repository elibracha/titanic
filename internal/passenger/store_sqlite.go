package passenger

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type sqliteStore struct {
	connector Connector
}

func (s *sqliteStore) GetPassenger(pid int) (*Passenger, error) {
	db, err := s.connector.Get()
	if err != nil {
		return nil, err
	}

	var passenger Passenger
	if err = db.Where("id = ?", pid).First(&passenger).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPassengerNotFound
		}
		return nil, err
	}
	return &passenger, nil
}

func (s *sqliteStore) GetPassengers() ([]*Passenger, error) {
	db, err := s.connector.Get()
	if err != nil {
		return nil, err
	}

	var passengers []*Passenger
	if err := db.Find(&passengers).Error; err != nil {
		// Handle the error
		fmt.Println("Error:", err.Error())
	}
	return passengers, nil
}

func NewStoreSQLite(connector Connector) Store {
	return &sqliteStore{connector: connector}
}
