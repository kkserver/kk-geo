package geo

import (
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

type Geo struct {
	IP          string  `json:"ip"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	TimeZone    string  `json:"timeZone"`
	Continent   string  `json:"continent"`
	CountryCode string  `json:"countryCode"`
	Country     string  `json:"country"`
	Province    string  `json:"province"`
	City        string  `json:"city"`
	PostalCode  string  `json:"postalCode"`
}

type GeoApp struct {
	app.App
	Geo         *GeoService
	Remote      *remote.Service
	Client      *client.Service
	ClientCache *client.WithService
	Expires     int64
	Geodb       string
}
