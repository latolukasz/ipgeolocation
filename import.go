package ipgeolocation

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

type ImportArguments struct {
	DbDirectory string
}

func Import(csvFilesDir string) error {

	if !strings.HasSuffix(csvFilesDir, "/") {
		csvFilesDir += "/"
	}
	err := importPlace(csvFilesDir)
	if err != nil {
		return err
	}
	err = importCountry(csvFilesDir)
	if err != nil {
		return err
	}
	err = importIpRange(csvFilesDir)
	if err != nil {
		return err
	}
	return nil
}

func importPlace(csvFilesDir string) error {
	rows, err := importFileAll(csvFilesDir, "db-place.csv")
	if err != nil {
		return err
	}
	places := make([]string, len(rows))
	for i, row := range rows {
		places[i] = row[1]
	}
	asJson, err := jsoniter.ConfigFastest.Marshal(places)
	if err != nil {
		return err
	}
	return saveFile(csvFilesDir, "db-place.db", asJson)
}

func convertStringToInt(val string) int {
	asInt, _ := strconv.Atoi(val)
	return asInt
}

func importCountry(csvFilesDir string) error {
	rows, err := importFileAll(csvFilesDir, "db-country.csv")
	if err != nil {
		return err
	}
	countries := make([][]interface{}, len(rows))
	for i, row := range rows {
		country := make([]interface{}, 6)
		country[0] = row[1]
		country[1] = convertStringToInt(row[2])
		country[2] = row[3]
		country[3] = convertStringToInt(row[5])
		country[4] = row[7]
		country[5] = strings.Trim(row[12], "")
		countries[i] = country
	}
	asJson, err := jsoniter.ConfigFastest.Marshal(countries)
	if err != nil {
		return err
	}
	return saveFile(csvFilesDir, "db-country.db", asJson)
}

func importIpRange(csvFilesDir string) error {
	f, err := os.Open(csvFilesDir + "db-ip-geolocation.csv")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	buff := bytes.Buffer{}
	timeZones := make(map[string]int)
	timeZoneID := 0

	reader := csv.NewReader(f)
	bs := make([]byte, 4)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ip := net.ParseIP(record[0]).To4()
		if ip == nil {
			break
		}
		_, _ = buff.Write(ip)
		ip = net.ParseIP(record[1]).To4()
		_, _ = buff.Write(ip)
		binary.LittleEndian.PutUint32(bs, uint32(convertStringToInt(record[2])))
		_, _ = buff.Write(bs)
		binary.LittleEndian.PutUint32(bs, uint32(convertStringToInt(record[5])))
		_, _ = buff.Write(bs)

		floatVal, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return err
		}
		lat := math.Float32bits(float32(floatVal))
		binary.LittleEndian.PutUint32(bs, lat)
		_, _ = buff.Write(bs)
		floatVal, err = strconv.ParseFloat(record[8], 64)
		if err != nil {
			return err
		}
		lon := math.Float32bits(float32(floatVal))
		binary.LittleEndian.PutUint32(bs, lon)
		_, _ = buff.Write(bs)

		timeZone := record[10]
		id, has := timeZones[timeZone]
		if has {
			binary.LittleEndian.PutUint32(bs, uint32(id))
			_, _ = buff.Write(bs)
		} else {
			binary.LittleEndian.PutUint32(bs, uint32(timeZoneID))
			_, _ = buff.Write(bs)
			timeZones[timeZone] = timeZoneID
			timeZoneID++
		}
	}
	timeZonesSlice := make([]string, len(timeZones))
	for k, v := range timeZones {
		timeZonesSlice[v] = k
	}
	asJson, err := jsoniter.ConfigFastest.Marshal(timeZonesSlice)
	if err != nil {
		return err
	}
	err = saveFile(csvFilesDir, "db-timezones.db", asJson)
	if err != nil {
		return err
	}
	return saveFile(csvFilesDir, "db-ip-geolocation.db", buff.Bytes())
}

func importFileAll(csvFilesDir string, fileName string) (records [][]string, err error) {
	f, err := os.Open(csvFilesDir + fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	return csv.NewReader(f).ReadAll()
}

func saveFile(csvFilesDir string, fileName string, data []byte) error {
	f, err := os.Create(csvFilesDir + fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write(data)
	return err
}
