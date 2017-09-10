package hksensibo

import (
	"github.com/llun/hksensibo/actions"
)

const (
	DRY_MODE  string = "dry"
	COOL_MODE        = "cool"
	AUTO_MODE        = "auto"
)

func (s *Sensibo) onTargetTemperatureUpdate(temperature float64) {
	action := actions.NewUpdateTargetTemperature(s.api, s.pod, s, int(temperature))
	s.worker.AddAction(action)
}

func (s *Sensibo) onHeatingCoolingStateUpdate(status int) {
	switch status {
	case 1:
		s.worker.AddAction(actions.NewUpdatePowerState(s.api, s.pod, s, true))
		s.worker.AddAction(actions.NewUpdateAcMode(s.api, s.pod, s, DRY_MODE))
	case 2:
		s.worker.AddAction(actions.NewUpdatePowerState(s.api, s.pod, s, true))
		s.worker.AddAction(actions.NewUpdateAcMode(s.api, s.pod, s, COOL_MODE))
	case 3:
		s.worker.AddAction(actions.NewUpdatePowerState(s.api, s.pod, s, true))
		s.worker.AddAction(actions.NewUpdateAcMode(s.api, s.pod, s, AUTO_MODE))
	default:
		s.worker.AddAction(actions.NewUpdatePowerState(s.api, s.pod, s, false))
	}
}
