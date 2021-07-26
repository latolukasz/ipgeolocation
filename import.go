package ipgeolocation

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/latolukasz/beeorm"
)

type ImportArguments struct {
	DbDirectory    string
	MysqlURI       string
	wrongCountryID uint64
	timeZones      map[string]uint16
	timeZoneID     uint16
}

func Import(ctx context.Context, arguments *ImportArguments) error {

	var placeEntity *PlaceEntity
	var countryEntity *CountryEntity
	var ipRangeV4Entity *IpRangeV4Entity
	var ipRangeV6Entity *IpRangeV6Entity
	var timeZoneEntity *TimeZoneEntity

	registry := beeorm.NewRegistry()
	registry.RegisterMySQLPool(arguments.MysqlURI)
	registry.RegisterEntity(placeEntity, countryEntity, ipRangeV4Entity, ipRangeV6Entity, timeZoneEntity)
	registry.RegisterEnumStruct("ipgeolocation.ContinentEnum", ContinentEnum, ContinentEnum.EU)

	validatedRegistry, err := registry.Validate(ctx)
	if err != nil {
		return err
	}
	engine := validatedRegistry.CreateEngine(ctx)

	validatedRegistry.GetTableSchemaForEntity(ipRangeV4Entity).DropTable(engine)
	validatedRegistry.GetTableSchemaForEntity(ipRangeV6Entity).DropTable(engine)
	validatedRegistry.GetTableSchemaForEntity(timeZoneEntity).DropTable(engine)
	validatedRegistry.GetTableSchemaForEntity(timeZoneEntity).UpdateSchema(engine)
	validatedRegistry.GetTableSchemaForEntity(ipRangeV4Entity).UpdateSchema(engine)
	validatedRegistry.GetTableSchemaForEntity(ipRangeV6Entity).UpdateSchema(engine)
	validatedRegistry.GetTableSchemaForEntity(countryEntity).UpdateSchemaAndTruncateTable(engine)
	validatedRegistry.GetTableSchemaForEntity(placeEntity).UpdateSchemaAndTruncateTable(engine)

	arguments.timeZones = make(map[string]uint16)
	flusher := engine.NewFlusher()
	if strings.HasSuffix(arguments.DbDirectory, "/") {
		arguments.DbDirectory += "/"
	}
	err = importPlace(flusher, arguments)
	if err != nil {
		return err
	}
	err = importCountry(flusher, arguments)
	if err != nil {
		return err
	}
	err = importIpRange(engine, flusher, arguments)
	if err != nil {
		return err
	}
	return nil
}

func importPlace(flusher beeorm.Flusher, arguments *ImportArguments) error {
	var placeEntity *PlaceEntity
	importer := func(record []string) (beeorm.Entity, error) {
		placeEntity = &PlaceEntity{}
		id, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		placeEntity.ID = uint32(id)
		placeEntity.NameEN = record[1]
		placeEntity.NameDE = record[2]
		placeEntity.NameRU = record[3]
		placeEntity.NameJA = record[4]
		placeEntity.NameFR = record[5]
		placeEntity.NameCN = record[6]
		placeEntity.NameES = record[7]
		placeEntity.NameCS = record[8]
		placeEntity.NameIT = record[9]
		return placeEntity, nil
	}
	return importFile(flusher, arguments, "db-place.csv", importer, 5000)
}

func importCountry(flusher beeorm.Flusher, arguments *ImportArguments) error {

	var countryEntity *CountryEntity
	importer := func(record []string) (beeorm.Entity, error) {
		id, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		if record[1] == "ZZ" {
			arguments.wrongCountryID = id
			return nil, nil
		}
		countryEntity = &CountryEntity{}

		countryEntity.ID = uint16(id)
		countryEntity.Continent = record[1]
		id, err = strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		countryEntity.ContinentName = &PlaceEntity{ID: uint32(id)}
		countryEntity.ISO2 = record[3]
		countryEntity.ISO3 = record[4]
		if record[5] != "" {
			id, err = strconv.ParseUint(record[5], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		countryEntity.Name = &PlaceEntity{ID: uint32(id)}
		if record[6] != "" {
			id, err = strconv.ParseUint(record[6], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		countryEntity.CapitalCityName = &PlaceEntity{ID: uint32(id)}
		countryEntity.CurrencyCode = record[7]
		countryEntity.CurrencyName = record[8]
		countryEntity.CurrencySuffix = record[9]
		countryEntity.CallingCode = record[10]
		countryEntity.Domain = record[11]
		countryEntity.Languages = record[12]
		return countryEntity, nil
	}
	return importFile(flusher, arguments, "db-country.csv", importer, 5000)
}

func importIpRange(engine *beeorm.Engine, flusher beeorm.Flusher, arguments *ImportArguments) error {
	var rangeEntity Range
	importer := func(record []string) (beeorm.Entity, error) {
		if record[2] == "" {
			return nil, nil
		}
		countryID, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return nil, err
		}
		if countryID == arguments.wrongCountryID {
			return nil, nil
		}
		rangeEntity = Range{}
		rangeEntity.Country = &CountryEntity{ID: uint16(countryID)}
		if record[3] != "" {
			refID, err := strconv.ParseUint(record[3], 10, 64)
			if err != nil {
				return nil, err
			}
			rangeEntity.State = &PlaceEntity{ID: uint32(refID)}
		}
		if record[4] != "" {
			refID, err := strconv.ParseUint(record[4], 10, 64)
			if err != nil {
				return nil, err
			}
			rangeEntity.District = &PlaceEntity{ID: uint32(refID)}
		}
		if record[5] != "" {
			refID, err := strconv.ParseUint(record[5], 10, 64)
			if err != nil {
				return nil, err
			}
			rangeEntity.City = &PlaceEntity{ID: uint32(refID)}
		}
		if record[7] != "" {
			floatVal, err := strconv.ParseFloat(record[7], 64)
			if err != nil {
				return nil, err
			}
			rangeEntity.Lat = float32(floatVal)
		}
		if record[8] != "" {
			floatVal, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				return nil, err
			}
			rangeEntity.Lon = float32(floatVal)
		}
		timeZoneID, has := arguments.timeZones[record[10]]
		if !has {
			arguments.timeZoneID++
			timeZoneID = arguments.timeZoneID
			timeZone := &TimeZoneEntity{Name: record[10], ID: arguments.timeZoneID}
			engine.Flush(timeZone)
			arguments.timeZones[record[10]] = timeZoneID
		}
		rangeEntity.TimeZone = &TimeZoneEntity{ID: timeZoneID}

		ip := net.ParseIP(record[0])
		if ip.To4() == nil {
			rangeV6Entity := &IpRangeV6Entity{Range: rangeEntity}
			rangeV6Entity.Start = ip2String(ip)
			ip = net.ParseIP(record[1])
			rangeV6Entity.End = ip2String(ip)
			return rangeV6Entity, nil
		}
		rangeV4Entity := &IpRangeV4Entity{Range: rangeEntity}
		rangeV4Entity.Start = ip2Uint(ip)
		ip = net.ParseIP(record[1])
		rangeV4Entity.End = ip2Uint(ip)
		return rangeV4Entity, nil
	}
	return importFile(flusher, arguments, "db-ip-geolocation.csv", importer, 5000)
}

type importRecord func(record []string) (beeorm.Entity, error)

func importFile(flusher beeorm.Flusher, arguments *ImportArguments, fileName string, importer importRecord, batch int) error {
	f, err := os.Open(arguments.DbDirectory + fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	reader := csv.NewReader(f)
	total := 0
	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		total++
		e, err := importer(record)
		if err != nil {
			return err
		}
		if e == nil {
			continue
		}
		flusher.Track(e)
		i++
		if i == batch {
			err = flusher.FlushInTransactionWithCheck()
			if err != nil {
				return err
			}
			flusher.Clear()
			i = 0
		}
	}
	err = flusher.FlushInTransactionWithCheck()
	if err != nil {
		return err
	}
	flusher.Clear()
	return nil
}

func ip2Uint(ip net.IP) uint32 {
	return func() uint32 {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("IP %v %v\n", ip, r)
			}
		}()
		ip = ip.To4()
		return binary.BigEndian.Uint32(ip)
	}()
}

func ip2String(ip net.IP) string {
	return func() string {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("IP2 %v %v\n", ip, r)
			}
		}()
		asString := strings.ReplaceAll(ip.To16().String(), "::", ":0000:00000:")
		parts := strings.Split(asString, ":")
		for i, part := range parts {
			if len(part) < 4 {
				parts[i] = strings.Repeat("0", 4-len(part)) + part
			}
		}
		asString = strings.Join(parts, ":")
		if len(parts) < 8 {
			asString += strings.Repeat(":0000", 8-len(parts))
		}
		return asString
	}()
}
