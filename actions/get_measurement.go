package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type GetMeasurement struct{}

func NewGetMeasurement() *GetMeasurement {
	return &GetMeasurement{}
}

func (a *GetMeasurement) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	measurements, err := api.GetMeasurements(pod.ID)
	if err != nil {
		log.Debug.Printf("Cannot get measurement for pod %v with error %v", pod.ID, err)
		return
	}

	if len(measurements) != 0 {
		store.UpdateMeasurement(measurements[0])
		log.Debug.Println("Update measurement to ", measurements[0])
	}
}

func (a *GetMeasurement) Name() string {
	return "GetMeasurement"
}

func (a *GetMeasurement) RemoveDuplicateCommand() bool {
	return false
}
