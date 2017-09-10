package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"time"
)

type UpdatePowerState struct {
	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store

	power bool
}

func NewUpdatePowerState(api *sensibo.Sensibo, pod sensibo.Pod, store Store, power bool) *UpdatePowerState {
	return &UpdatePowerState{api, pod, store, power}
}

func (a *UpdatePowerState) Run() {
	state := a.store.CurrentAcState()
	state.On = a.power
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

func (a *UpdatePowerState) Name() string {
	return "UpdatePowerState"
}

func (a *UpdatePowerState) RemoveDuplicateCommand() bool {
	return true
}
