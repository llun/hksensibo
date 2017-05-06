package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type UpdateAcState struct {
	state sensibo.AcState
}

func NewUpdateAcState(state sensibo.AcState) *UpdateAcState {
	return &UpdateAcState{state}
}

func (a *UpdateAcState) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	log.Debug.Printf("Update %v to %v", pod.ID, a.state)
	store.UpdateAcState(a.state)
	api.ReplaceState(pod.ID, a.state)
}

func (a *UpdateAcState) Name() string {
	return "UpdateAcState"
}

func (a *UpdateAcState) RemoveDuplicateCommand() bool {
	return false
}
