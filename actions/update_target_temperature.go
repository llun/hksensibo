package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"time"
)

type UpdateTargetTemperature struct {
	temperature int
}

func NewUpdateTargetTemperature(temperature int) *UpdateTargetTemperature {
	return &UpdateTargetTemperature{temperature}
}

func (a *UpdateTargetTemperature) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	state := store.CurrentAcState()
	state.TargetTemperature = a.temperature
	store.UpdateAcState(state)

	log.Debug.Printf("Sensibo Update %v to %v", pod.ID, state)
	for i := 0; i < RETRY_COUNT; i++ {
		response, err := api.ReplaceState(pod.ID, state)
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
