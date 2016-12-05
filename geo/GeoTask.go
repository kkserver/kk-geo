package geo

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type GeoTaskResult struct {
	app.Result
	Geo *Geo `json:"geo,omitempty"`
}

type GeoTask struct {
	app.Task
	IP     string `json:"ip"`
	Result GeoTaskResult
}

func (task *GeoTask) GetResult() interface{} {
	return &task.Result
}

func (task *GeoTask) GetInhertType() string {
	return "geo"
}

func (task *GeoTask) GetClientName() string {
	return "Geo.Get"
}
