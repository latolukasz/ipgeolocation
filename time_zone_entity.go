package ipgeolocation

import "github.com/latolukasz/beeorm"

type TimeZoneEntity struct {
	beeorm.ORM
	ID   uint16
	Name string `orm:"required"`
}
