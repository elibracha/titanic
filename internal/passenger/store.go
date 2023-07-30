package passenger

const (
	CSV    StoreType = "CSV"
	SQLITE StoreType = "SQLITE"
)

type StoreType string

type Passenger struct {
	PassengerId int     `csv:"PassengerId"gorm:"column:PassengerId;primary_key"`
	Survived    int     `csv:"Survived"gorm:"column:Survived"`
	Pclass      int     `csv:"Pclass"gorm:"column:Pclass"`
	Name        string  `csv:"Name"gorm:"column:Name"`
	Sex         string  `csv:"Sex"gorm:"column:Sex"`
	Age         string  `csv:"Age"gorm:"column:Age"`
	SibSp       int     `csv:"SibSp"gorm:"column:SibSp"`
	Parch       int     `csv:"Parch"gorm:"column:Parch"`
	Ticket      string  `csv:"Ticket"gorm:"column:Ticket"`
	Fare        float64 `csv:"Fare"gorm:"column:Fare"`
	Cabin       string  `csv:"Cabin"gorm:"column:Cabin"`
	Embarked    string  `csv:"Embarked"gorm:"column:Embarked"`
}

type Store interface {
	GetPassengers() ([]*Passenger, error)
	GetPassenger(pid int) (*Passenger, error)
}
