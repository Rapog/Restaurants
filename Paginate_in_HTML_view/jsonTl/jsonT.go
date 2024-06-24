package jsonTl

import (
	"strconv"
)

type Place struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location Locs   `json:"location"`
}

type Locs struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

func FromCSVtoStruct(record []string, place *Place) {
	place.Id, _ = strconv.ParseInt(record[0], 10, 64)
	place.Name = record[1]
	place.Address = record[2]
	place.Phone = record[3]
	place.Location.Longitude, _ = strconv.ParseFloat(record[4], 64)
	place.Location.Latitude, _ = strconv.ParseFloat(record[5], 64)
}
