package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"time"
)

type UpdateTargetTemperature struct {
	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store

	temperature int
}

func NewUpdateTargetTemperature(api *sensibo.Sensibo, pod sensibo.Pod, store Store, temperature int) *UpdateTargetTemperature {
	return &UpdateTargetTemperature{api, pod, store, temperature}
}

func (a *UpdateTargetTemperature) Run() {
	state := a.store.CurrentAcState()
	state.TargetTemperature = a.temperature
	a.store.UpdateAcState(state)

	log.Debug.Printf("Sensibo Update %v to %v", a.pod.ID, state)
	for i := 0; i < RETRY_COUNT; i++ {
		response, err := a.api.ReplaceState(a.pod.ID, state)
		log.Debug.Println("Sensibo response", response)
		if err != nil {
			log.Debug.Println("Sensibo error", err)

			// Don't retry immediatly
			wait := make(chan bool)
			time.AfterFunc(1*time.Second, func() { wait <- true })
			<-wait
		} else {
			break
		}
	}
}

func (a *UpdateTargetTemperature) Name() string {
	return "UpdateTargetTemperature"
}

func (a *UpdateTargetTemperature) RemoveDuplicateCommand() bool {
	return true
}
