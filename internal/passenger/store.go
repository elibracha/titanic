package passenger

import "fmt"

const (
	CSV    StoreType = "CSV"
	SQLITE StoreType = "SQLITE"
)

var (
	ErrPassengerNotFound = fmt.Errorf("passenger not found")
)

type StoreType string

type Passenger struct {
	PassengerId int     `csv:"PassengerId"gorm:"column:id;primary_key"`
	Survived    int     `csv:"Survived"gorm:"column:survived"`
	Pclass      int     `csv:"Pclass"gorm:"column:class"`
	Name        string  `csv:"Name"gorm:"column:name"`
	Sex         string  `csv:"Sex"gorm:"column:sex"`
	Age         string  `csv:"Age"gorm:"column:age"`
	SibSp       int     `csv:"SibSp"gorm:"column:siblings_spouses"`
	Parch       int     `csv:"Parch"gorm:"column:parents_children"`
	Ticket      string  `csv:"Ticket"gorm:"column:ticket"`
	Fare        float64 `csv:"Fare"gorm:"column:fare"`
	Cabin       string  `csv:"Cabin"gorm:"column:cabin"`
	Embarked    string  `csv:"Embarked"gorm:"column:embarked"`
}

type Store interface {
	GetPassengers() ([]*Passenger, error)
	GetPassenger(pid int) (*Passenger, error)
}
