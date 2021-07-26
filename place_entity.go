package ipgeolocation

import "github.com/latolukasz/beeorm"

type PlaceEntity struct {
	beeorm.ORM
	ID     uint32
	NameEN string `orm:"required"`
	NameDE string `orm:"required"`
	NameRU string `orm:"required"`
	NameJA string `orm:"required"`
	NameFR string `orm:"required"`
	NameCN string `orm:"required"`
	NameES string `orm:"required"`
	NameCS string `orm:"required"`
	NameIT string `orm:"required"`
}
