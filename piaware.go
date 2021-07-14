package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slim-bean/adsb-loki/pkg/aircraft"
	"github.com/slim-bean/adsb-loki/pkg/model"
	"github.com/slim-bean/adsb-loki/pkg/piaware"
	"os"
	"telegraf_piaware/internal/aircraft_metric"
	icao86432 "telegraf_piaware/internal/icao8643"
	"telegraf_piaware/internal/tclogger"
	"time"
)
var telegrafLogger = tclogger.Create().Start()
var localLog log.Logger
var icao8643 = icao86432.New()

func main() {
	var piawareUrl string
	var aircraftUrl string


	localLog = log.NewLogfmtLogger(log.NewSyncWriter(telegrafLogger))
	localLog = log.With(localLog, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	flag.StringVar(&piawareUrl, "url", "http://localhost/dump1090/data/aircraft.json", "URL pointing to aircraft.json on piaware system")
	flag.StringVar(&aircraftUrl, "aircraft", "https://github.com/wiedehopf/tar1090-db/raw/csv/aircraft.csv.gz", "URL to the aircraft database (.csv.gz)")

	var config = aircraft.Config{
		Directory:  "c:/work/go/telegraf-piaware",
		BoltDbFile: "c:/work/go/telegraf-piaware/db",
		URL: aircraftUrl,
	}

	flag.Parse()



	/*
	 * The aircraft manager gets the initial database, then starts a thread
	 * that periodically fetches it again to keep it up to date
	 */
	m, err := aircraft.NewAircraftManager(localLog, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init the registration manager: %v\n", err)
		os.Exit(1)
	}
	m.Run()

	//reg, err := registration.NewManager(localLog, registration.RegManagerConfig{
	//	Directory: ".",
	//	URL:       "http://registry.faa.gov/database/ReleasableAircraft.zip",
	//})




	err = icao8643.Load()
	if err != nil {
		level.Error(localLog).Log("msg", "Unable to fetch ICAO 8643 aircraft type database")
	}

	p := piaware.New(m,piawareUrl )

	getData(p)
}

func getData(p *piaware.Piaware) {

	report, err := p.GetReport()
	if err != nil {
		panic(err)
	}


	if report == nil {
		level.Info(localLog).Log("msg", "Report was empty")
		return
}
		reportTime := time.Unix(int64(report.Now),0)

		var metrics aircraft_metric.AircraftMetricList

		for _, a := range report.Aircraft {
			metric  := aircraftObservationToMetric(a, reportTime)
			metrics = append(metrics, metric)
		}

		metrics.WriteTo(os.Stdout)
}

func aircraftObservationToMetric(aircraft model.Aircraft, timestamp time.Time) aircraft_metric.AircraftMetric {

	am := aircraft_metric.New("adsb", timestamp)
	am.AddTag("hex", &aircraft.Hex)
	am.AddTag("type", aircraft.TypeCode)
	am.AddTag("description", aircraft.Description)
	am.AddTag("category", aircraft.Category)
	am.AddTag("registration", aircraft.Registration)
	am.AddBoolTagIfTrue("military", aircraft.Military)
	am.AddBoolTagIfTrue("interesting", aircraft.Interesting)
	am.AddBoolTagIfTrue("pia", aircraft.PIA)
	am.AddBoolTagIfTrue("ladd", aircraft.LADD)

	if icao8643 != nil && aircraft.TypeCode != nil {
		idata:= icao8643.Get(*aircraft.TypeCode)
		if idata != nil {
			am.AddTag("icao_description",idata.AircraftDescription)
			am.AddTag("icao_manufacturer", idata.ManufacturerCode)
			am.AddTag("icao_model", idata.ModelFullName)
			am.AddTag("icao_engine_count", idata.EngineCount)
			am.AddTag("icao_engine_type", idata.EngineType)
		}
	}


	am.AddField("lat", aircraft.Lat)
	am.AddField("lng", aircraft.Lon)
	am.AddField("track", aircraft.Track)
	am.AddField("emergency", aircraft.Emergency)
	am.AddField("alt", aircraft.BarometerAltitude)
	am.AddField("flight", aircraft.Flight)
	am.AddField("geom_altitude", aircraft.GeometricAltitude)
	am.AddField("squawk", aircraft.Squawk)
	am.AddField("rssi", aircraft.Rssi)
	am.AddField("ground_speed", aircraft.GroundSpeed)





	return am
}

