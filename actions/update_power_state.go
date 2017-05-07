package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"time"
)

type UpdatePowerState struct {
	power bool
}

func NewUpdatePowerState(power bool) *UpdatePowerState {
	return &UpdatePowerState{power}
}

func (a *UpdatePowerState) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	state := store.CurrentAcState()
	state.On = a.power
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

func (a *UpdatePowerState) Name() string {
	return "UpdatePowerState"
}

func (a *UpdatePowerState) RemoveDuplicateCommand() bool {
	return true
}
