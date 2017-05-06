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
	action := actions.NewUpdateTargetTemperature(int(temperature))
	s.worker.AddAction(action)
}

func (s *Sensibo) onHeatingCoolingStateUpdate(status int) {
	switch status {
	case 1:
		s.worker.AddAction(actions.NewUpdatePowerState(true))
		s.worker.AddAction(actions.NewUpdateAcMode(DRY_MODE))
	case 2:
		s.worker.AddAction(actions.NewUpdatePowerState(true))
		s.worker.AddAction(actions.NewUpdateAcMode(COOL_MODE))
	case 3:
		s.worker.AddAction(actions.NewUpdatePowerState(true))
		s.worker.AddAction(actions.NewUpdateAcMode(AUTO_MODE))
	default:
		s.worker.AddAction(actions.NewUpdatePowerState(false))
	}
}
