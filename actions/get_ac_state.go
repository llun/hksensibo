package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

type GetAcState struct{}

func NewGetAcState() *GetAcState {
	return &GetAcState{}
}

func (a *GetAcState) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	states, err := api.GetAcStates(pod.ID)
	if err != nil {
		log.Debug.Printf("Cannot get ac state for pod %v with error %v", pod.ID, err)
		return
	}

	if len(states) != 0 {
		store.UpdateAcState(states[0].AcState)
		log.Debug.Println("Update state to ", states[0].AcState)
	}
}

func (a *GetAcState) Name() string {
	return "GetAcState"
}

func (a *GetAcState) RemoveDuplicateCommand() bool {
	return false
}
