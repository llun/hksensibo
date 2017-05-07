package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
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

	log.Debug.Printf("Update %v to %v", pod.ID, state)
	response, err := api.ReplaceState(pod.ID, state)
	if err != nil {
		log.Debug.Println("Sensibo error", err)
	}
	log.Debug.Println("Sensibo response", response)
}

func (a *UpdatePowerState) Name() string {
	return "UpdatePowerState"
}

func (a *UpdatePowerState) RemoveDuplicateCommand() bool {
	return true
}
