package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"time"
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

func (a *UpdateAcMode) Name() string {
	return "UpdateAcMode"
}

func (a *UpdateAcMode) RemoveDuplicateCommand() bool {
	return true
}
