package ipgeolocation

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
)

var dbIP []byte
var countries [][]interface{}
var names []string
var timezones []string

func InitDB(dbDirectory string) error {

	if !strings.HasSuffix(dbDirectory, "/") {
		dbDirectory += "/"
	}
	var err error
	dbIP, err = os.ReadFile(dbDirectory + "db-ip-geolocation.db")
	if err != nil {
		return err
	}
	countriesData, err := os.ReadFile(dbDirectory + "db-country.db")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigFastest.Unmarshal(countriesData, &countries)
	if err != nil {
		return err
	}
	namesData, err := os.ReadFile(dbDirectory + "db-place.db")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigFastest.Unmarshal(namesData, &names)
	if err != nil {
		return err
	}
	timezonesData, err := os.ReadFile(dbDirectory + "db-timezones.db")
	if err != nil {
		return err
	}
	err = jsoniter.ConfigFastest.Unmarshal(timezonesData, &timezones)
	if err != nil {
		return err
	}
	return nil
}
