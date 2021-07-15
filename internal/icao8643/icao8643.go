package icao8643

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Icao8643 struct {
	byTypeCode map[string]AircraftType

}

func New() *Icao8643 {
	return &Icao8643{
		byTypeCode: make(map[string]AircraftType),
	}
}


func (i *Icao8643) Load() error {
	response, err := http.Post("https://www4.icao.int/doc8643/External/AircraftTypes", "text/json", nil)

	if err != nil {
		return err
	}
	responseBytes, err:= ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var aircraftTypes aircraftTypeList
	err = json.Unmarshal(responseBytes, &aircraftTypes)
	if err != nil {
		return err
	}

	for _, a := range aircraftTypes {
		i.byTypeCode[a.Designator] = a
	}

	return nil
}


func (i *Icao8643) Get(aircraftTypeCode string) *AircraftType {

	a, exists := i.byTypeCode[aircraftTypeCode]

	if exists {
		return &a
	} else {
		return nil
	}

}