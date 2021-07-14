package aircraft_metric

import (
	"bytes"
	protocol "github.com/influxdata/line-protocol"
	"io"
	"time"
)

type AircraftMetric struct {
	timestamp time.Time
	name      string
	tags      map[string]string
	fields    map[string]interface{}
}
type AircraftMetricList []AircraftMetric


func New(name string, reportTime time.Time) AircraftMetric {
	return AircraftMetric{
		timestamp: reportTime,
		name:      name,
		tags:      make(map[string]string),
		fields:    make(map[string]interface{}),
	}
}

func (a *AircraftMetric) AddTag(key string, value *string) {
	if value != nil {
		a.tags[key] = *value
	}
}


func (a *AircraftMetric) AddBoolTagIfTrue(key string, value *bool) {
	if value != nil {
		if *value {
			a.tags[key] = "y"
		}
	}
}

func (a *AircraftMetric) AddField(key string, value interface{}) {
	if value != nil {
		a.fields[key] = value
	}
}


func (a *AircraftMetric) CreateMetric() (protocol.Metric, error) {
	return protocol.New(a.name, a.tags, a.fields, a.timestamp)
}

func (a *AircraftMetric) WriteTo(out io.Writer) (int, error) {
	buf := &bytes.Buffer{}
	serializer := protocol.NewEncoder(buf)
	serializer.SetMaxLineBytes(-1)
	serializer.SetFieldTypeSupport(protocol.UintSupport)

	m, err := a.CreateMetric()
	if err != nil {
		return 0, err
	}

	_, err = serializer.Encode(m)
	if err != nil {
		return 0, err
	}

	return out.Write(buf.Bytes())
}



func (metrics AircraftMetricList) WriteTo(out io.Writer) (int, error) {
	buf := &bytes.Buffer{}
	serializer := protocol.NewEncoder(buf)
	serializer.SetMaxLineBytes(-1) /* let it get as big as it needs to */
	serializer.SetFieldTypeSupport(protocol.UintSupport)

	for _, metric := range metrics {
		im, err := metric.CreateMetric()
		if err != nil {
			continue
		}

		_, err = serializer.Encode(im)
		if err != nil {
			continue
		}

	}
	return out.Write(buf.Bytes())
}

