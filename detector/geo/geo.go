package geo

import (
	"github.com/oschwald/geoip2-golang"
)

func NewGeo(dataSourceName string) (*geoip2.Reader, error) {
	db, err := geoip2.Open(dataSourceName)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}