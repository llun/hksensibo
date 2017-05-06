package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"
)

const (
	DRY_MODE  string = "dry"
	COOL_MODE        = "cool"
	AUTO_MODE        = "auto"
)

type UpdateAcMode struct {
	mode string
}

func NewUpdateAcMode(mode string) *UpdateAcMode {
	return &UpdateAcMode{mode}
}

func (a *UpdateAcMode) Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store) {
	state := store.CurrentAcState()
	state.Mode = a.mode
	store.UpdateAcState(state)

	log.Debug.Printf("Update %v to %v", pod.ID, state)
	api.ReplaceState(pod.ID, state)
}

func (a *UpdateAcMode) Name() string {
	return "UpdateAcMode"
}

func (a *UpdateAcMode) RemoveDuplicateCommand() bool {
	return true
}
