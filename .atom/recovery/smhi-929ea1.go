package smhi

import (
	"bytes"
	"time"
	"encoding/binary"
)

const URL string = "https://opendata-download-metfcst.smhi.se/api/category/pmp3g/version/2/geotype/point/lon/%v/lat/%v/data.json"

type Entry struct {
	City string
	ReferencedTime time.Time
	Longitude float64
	Latitude float64
}

func (e Entry) Store() (int64, error) {
  res, err := ForecastDB.Exec(`
    INSERT INTO
      forecast.Smhi(reference_time, city, longitude, latitude)
    VALUES
      (?, ?, ?, ?);`,
		e.ReferenceTime,
		e.City,
		e.Longitude,
		e.Latitude)
  if err != nil {
    return err
  }
	return res.LastInsertId()
}

type Row struct {
	TimeMeasured   float64
	TimeReferenced float64
	Pressure       float64
	Temperature    float64
	Visibility     float64
	WindSpeed      float64
	WindDirection  float64
	Humidity       float64
	CloudCoverage  float64
}


func (r Row) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, r); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
// https://stackoverflow.com/questions/34701187/go-byte-to-little-big-endian-signed-integer-or-float
func (r Row) Decode(data []byte) error {
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &r); err != nil {
		return err
	}
	return nil
}
