package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type GetAcState struct {
	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store
}

func NewGetAcState(api *sensibo.Sensibo, pod sensibo.Pod, store Store) *GetAcState {
	return &GetAcState{api, pod, store}
}

func (a *GetAcState) Run() {
	states, err := a.api.GetAcStates(a.pod.ID)
	if err != nil {
		log.Debug.Printf("Cannot get ac state for pod %v with error %v", a.pod.ID, err)
		return
	}

	if len(states) != 0 {
		a.store.UpdateAcState(states[0].AcState)
		log.Debug.Println("Update state to ", states[0].AcState)
	}
}

func (a *GetAcState) Name() string {
	return "GetAcState"
}

func (a *GetAcState) RemoveDuplicateCommand() bool {
	return false
}
