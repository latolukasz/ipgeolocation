package ipgeolocation

import "github.com/latolukasz/beeorm"

type CountryEntity struct {
	beeorm.ORM
	ID        uint16
	Continent string `orm:"required;length=2"`
}
