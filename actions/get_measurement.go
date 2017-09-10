package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type GetMeasurement struct {
	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store
}

func NewGetMeasurement(api *sensibo.Sensibo, pod sensibo.Pod, store Store) *GetMeasurement {
	return &GetMeasurement{api, pod, store}
}

func (a *GetMeasurement) Run() {
	measurements, err := a.api.GetMeasurements(a.pod.ID)
	if err != nil {
		log.Debug.Printf("Cannot get measurement for pod %v with error %v", a.pod.ID, err)
		return
	}

	if len(measurements) != 0 {
		a.store.UpdateMeasurement(measurements[0])
		log.Debug.Println("Update measurement to ", measurements[0])
	}
}

func (a *GetMeasurement) Name() string {
	return "GetMeasurement"
}

func (a *GetMeasurement) RemoveDuplicateCommand() bool {
	return false
}
