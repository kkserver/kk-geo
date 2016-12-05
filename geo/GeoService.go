package geo

import (
	"fmt"
	"github.com/kkserver/kk-cache/cache"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/json"
	geoip2 "github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

type GeoService struct {
	app.Service
	Init     *app.InitTask
	Get      *GeoTask
	geodb    *geoip2.Reader
	dispatch *kk.Dispatch
}

func (S *GeoService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *GeoService) HandleInitTask(a *GeoApp, task *app.InitTask) error {

	db, err := geoip2.Open(a.Geodb)

	if err != nil {
		log.Println("[GeoService][HandleInitTask]" + err.Error())
		return nil
	}

	S.geodb = db
	S.dispatch = kk.NewDispatch()

	return nil
}

func (S *GeoService) HandleGeoTask(a *GeoApp, task *GeoTask) error {

	if task.IP == "" {
		task.Result.Errno = ERROR_GEO_NOT_FOUND_IP
		task.Result.Errmsg = "Not found ip"
		return nil
	}

	{
		var cache = cache.CacheTask{}
		cache.Key = fmt.Sprintf("geo.%s", task.IP)
		app.Handle(a, &cache)
		if cache.Result.Errno == 0 && cache.Result.Value != "" {
			var v = Geo{}
			err := json.Decode([]byte(cache.Result.Value), &v)
			if err == nil {
				task.Result.Geo = &v
				return nil
			}
		}
	}

	if S.dispatch != nil {

		S.dispatch.Sync(func() {

			var v = Geo{}

			v.IP = task.IP

			ip := net.ParseIP(v.IP)

			city, err := S.geodb.City(ip)
			if err != nil {
				task.Result.Errno = ERROR_GEO
				task.Result.Errmsg = err.Error()
				return
			} else {
				v.Continent = city.Continent.Names["en"]
				v.Country = city.Country.Names["en"]
				v.CountryCode = city.Country.IsoCode
				if len(city.Subdivisions) > 0 {
					v.Province = city.Subdivisions[0].Names["en"]
				}
				v.City = city.City.Names["en"]
				v.PostalCode = city.Postal.Code
				v.TimeZone = city.Location.TimeZone
				v.Latitude = city.Location.Latitude
				v.Longitude = city.Location.Longitude
				task.Result.Geo = &v
				return
			}

		})

		if task.Result.Geo != nil {
			var cache = cache.CacheSetTask{}
			cache.Key = fmt.Sprintf("geo.%s", task.IP)
			cache.Expires = a.Expires
			b, _ := json.Encode(task.Result.Geo)
			cache.Value = string(b)
			app.Handle(a, &cache)
		}
	}

	return nil
}
