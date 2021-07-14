package icao8643

type AircraftType struct {
	Designator          string
	ModelFullName       *string
	Description         *string
	WTC                 *string
	WTG                 *string
	ManufacturerCode    *string
	AircraftDescription *string
	EngineCount         *string
	EngineType          *string
}

type aircraftTypeList []AircraftType
