package smhi

import (
	"time"
	"fmt"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/remnestal/forelcastro/utils"
	"github.com/remnestal/forelcastro/data"
)

const URL string = "https://opendata-download-metfcst.smhi.se/api/category/pmp3g/version/2/geotype/point/lon/%v/lat/%v/data.json"

type Source struct {
	Log *logrus.Entry
	Pause time.Duration
}

func (s Source) Logger() *logrus.Entry {
	return s.Log
}

func (s Source) Interval() time.Duration {
	return s.Pause
}

type Param struct {
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Level     int       `json:"level"`
	LevelType string    `json:"levelType"`
	Values    []float64 `json:"values"`
}

type TimeSeries struct {
	ValidTime  time.Time `json:"validTime"`
	Parameters []Param   `json:"parameters"`
}

type Geodata struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Forecast struct {
	ApprovedTime  time.Time    `json:"approvedTime"`
	ReferenceTime time.Time    `json:"referenceTime"`
	Geometry      Geodata      `json:"geometry"`
	TimeSeries    []TimeSeries `json:"timeSeries"`
}

func float64ptr(f float64) *float64 {
	return &f
}

func (s Source) Ingest() {
	for _, gp := range data.GeoPoints {
		logger := s.Log.WithField("point", gp)
		payload := utils.Get(fmt.Sprintf(URL, gp.Lon, gp.Lat), logger)
		if payload == nil {
			logger.Errorln("Unable to fetch data")
			continue
		}
		var forecast Forecast
		if err := json.Unmarshal(payload, &forecast); err != nil {
			logger.WithError(err).Errorln("Unable to parse data")
			continue
		}
		id, err := Entry{
			City: gp.Label,
			ReferencedTime: forecast.ReferenceTime,
			Longitude: gp.Lon,
			Latitude: gp.Lat,
		}.Store()
		if err != nil {
			logger.WithError(err).Errorln("Unable to write entry in database")
			continue
		}
		for _, ts := range forecast.TimeSeries {
			r := Row{ ValidTime: ts.ValidTime }
			r.Populate(ts.Params)
			if _, err := r.Store(id); err != nil {
				s.Log.WithError(err).Errorln("Unable to write point in database")
				continue
			}
		}
	}

}
