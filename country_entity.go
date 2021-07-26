package ipgeolocation

import "github.com/latolukasz/beeorm"

var ContinentEnum = struct {
	EU string
	AS string
	NA string
	AF string
	AN string
	SA string
	OC string
	ZZ string
}{
	EU: "EU",
	AS: "AS",
	NA: "NA",
	AF: "AF",
	AN: "AN",
	SA: "SA",
	OC: "OC",
	ZZ: "ZZ",
}

type CountryEntity struct {
	beeorm.ORM
	ID              uint16
	Continent       string       `orm:"enum=ipgeolocation.ContinentEnum;required"`
	ISO2            string       `orm:"required;length=2"`
	ISO3            string       `orm:"required;length=3"`
	ContinentName   *PlaceEntity `orm:"required"`
	Name            *PlaceEntity `orm:"required"`
	CapitalCityName *PlaceEntity `orm:"required"`
	CurrencyCode    string       `orm:"required;length=3"`
	CurrencyName    string       `orm:"required"`
	CurrencySuffix  string       `orm:"required"`
	CallingCode     string       `orm:"required"`
	Domain          string
	Languages       string `orm:"required"`
}
