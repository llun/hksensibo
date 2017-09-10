package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type UpdateAcState struct {
	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store

	state sensibo.AcState
}

func NewUpdateAcState(api *sensibo.Sensibo, pod sensibo.Pod, store Store, state sensibo.AcState) *UpdateAcState {
	return &UpdateAcState{api, pod, store, state}
}

func (a *UpdateAcState) Run() {
	log.Debug.Printf("Update %v to %v", a.pod.ID, a.state)
	a.store.UpdateAcState(a.state)
	a.api.ReplaceState(a.pod.ID, a.state)
}

func (a *UpdateAcState) Name() string {
	return "UpdateAcState"
}

func (a *UpdateAcState) RemoveDuplicateCommand() bool {
	return false
}
