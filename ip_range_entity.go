package ipgeolocation

import "github.com/latolukasz/beeorm"

type Range struct {
	Country  *CountryEntity `orm:"required"`
	State    *PlaceEntity
	District *PlaceEntity
	City     *PlaceEntity
	Lat      float32         `orm:"decimal=9,5;unsigned=false"`
	Lon      float32         `orm:"decimal=9,5;unsigned=false"`
	TimeZone *TimeZoneEntity `orm:"required"`
}

type IpRangeV4Entity struct {
	beeorm.ORM
	ID    uint32
	Start uint32 `orm:"required"`
	End   uint32 `orm:"required"`
	Range
}

type IpRangeV6Entity struct {
	beeorm.ORM
	ID    uint32
	Start string `orm:"required;length=100"`
	End   string `orm:"required;length=100"`
	Range
}
