package sensibo

import (
  "github.com/brutella/hc/characteristic"
)

const (
  DRY_MODE  string = "dry"
  COOL_MODE        = "cool"
  AUTO_MODE        = "auto"
)

func (s *Sensibo) updateCurrentHeatingCoolingState() {
  state := s.Thermostat.CurrentHeatingCoolingState
  currentValue := characteristic.CurrentHeatingCoolingStateOff
  if s.CurrentState.On {
    switch s.CurrentState.Mode {
    case DRY_MODE:
      currentValue = characteristic.CurrentHeatingCoolingStateHeat
    default:
      currentValue = characteristic.CurrentHeatingCoolingStateCool
    }
  }
  state.SetValue(currentValue)
}

func (s *Sensibo) updateTargetHeatingCoolingState() {
  state := s.Thermostat.TargetHeatingCoolingState
  currentValue := characteristic.TargetHeatingCoolingStateOff
  if s.CurrentState.On {
    switch s.CurrentState.Mode {
    case DRY_MODE:
      currentValue = characteristic.TargetHeatingCoolingStateHeat
    case COOL_MODE:
      currentValue = characteristic.TargetHeatingCoolingStateCool
    default:
      currentValue = characteristic.TargetHeatingCoolingStateAuto
    }
  }
  state.SetValue(currentValue)
}

func (s *Sensibo) setupTargetHeatingCoolingState() {
  s.updateTargetHeatingCoolingState()
  state := s.Thermostat.TargetHeatingCoolingState
  state.OnValueRemoteUpdate(func(status int) {
    currentState := s.CurrentState
    switch status {
    case 1:
      currentState.On = true
      currentState.Mode = DRY_MODE
    case 2:
      currentState.On = true
      currentState.Mode = COOL_MODE
    case 3:
      currentState.On = true
      currentState.Mode = AUTO_MODE
    default:
      currentState.On = false
    }
    s.api.ReplaceState(s.pod.ID, currentState)
  })
}
