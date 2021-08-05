package ipgeolocation

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"strings"
)

type Result struct {
	Country  Country
	City     string
	Lat      float32
	Lon      float32
	Timezone string
}

type Continent struct {
	ISO  string
	Name string
}

type Country struct {
	Continent Continent
	ISO       string
	Name      string
	Currency  string
	Languages []string
}

func Search(ip string) (*Result, error) {
	netIP := net.ParseIP(ip).To4()
	if netIP == nil {
		return nil, nil
	}

	l := len(dbIP)
	page := 0
	for {
		start := page
		if bytes.Compare(netIP, dbIP[start:start+4]) >= 0 && bytes.Compare(netIP, dbIP[start+4:start+8]) <= 0 {
			res := &Result{}
			res.City = names[binary.LittleEndian.Uint32(dbIP[start+12:start+16])-1]
			country := countries[binary.LittleEndian.Uint32(dbIP[start+8:start+12])-1]
			res.Lat = math.Float32frombits(binary.LittleEndian.Uint32(dbIP[start+16 : start+20]))
			res.Lon = math.Float32frombits(binary.LittleEndian.Uint32(dbIP[start+20 : start+24]))
			res.Timezone = timezones[binary.LittleEndian.Uint32(dbIP[start+24:start+28])]
			res.Country = Country{
				Continent: Continent{
					ISO:  country[0].(string),
					Name: names[int(country[1].(float64))-1],
				},
				ISO:       country[2].(string),
				Name:      names[int(country[3].(float64))-1],
				Currency:  country[4].(string),
				Languages: strings.Split(country[5].(string), ","),
			}
			return res, nil
		}
		page += 28
		if page == l {
			break
		}
	}

	return nil, nil
}
