package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
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

	log.Debug.Printf("Update %v to %v", pod.ID, state)
	api.ReplaceState(pod.ID, state)
}

func (a *UpdateTargetTemperature) Name() string {
	return "UpdateTargetTemperature"
}

func (a *UpdateTargetTemperature) RemoveDuplicateCommand() bool {
	return true
}