package sensibo

import (
	"github.com/brutella/hc/characteristic"
)

func CelciusToFahrenheit(celcius float64) float64 {
	return celcius*9/5 + 32
}

func FahrenheitToCelcius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

func (s *Sensibo) setupTemperatureDisplayUnits() {
	s.updateTemperatureDisplayUnits()
	state := s.Thermostat.TemperatureDisplayUnits
	state.OnValueRemoteUpdate(func(status int) {
		if status == 0 {
			s.TemperatureUnit = TemperatureUnitCelcius
		} else {
			s.TemperatureUnit = TemperatureUnitFahrenheit
		}
		s.updateHomeKitFromState()
	})
}

func (s *Sensibo) setupTargetTemperature() {
	s.updateTargetTemperature()
	state := s.Thermostat.TargetTemperature
	state.OnValueRemoteUpdate(func(status float64) {
		newValue := status
		if s.TemperatureUnit == TemperatureUnitFahrenheit {
			newValue = FahrenheitToCelcius(newValue)
		}

		currentState := s.CurrentState
		currentState.TargetTemperature = int(newValue)
		go func() {
			s.api.ReplaceState(s.pod.ID, currentState)
		}()
	})
}

func (s *Sensibo) updateTemperatureDisplayUnits() {
	state := s.Thermostat.TemperatureDisplayUnits
	value := characteristic.TemperatureDisplayUnitsCelsius
	if s.TemperatureUnit == TemperatureUnitFahrenheit {
		value = characteristic.TemperatureDisplayUnitsFahrenheit
	}
	state.SetValue(value)
}

func (s *Sensibo) updateCurrentTemperature() {
	state := s.Thermostat.CurrentTemperature
	value := s.CurrentMeasurement.Temperature
	if s.TemperatureUnit == TemperatureUnitFahrenheit {
		value = CelciusToFahrenheit(value)
	}
	state.SetValue(value)
}

func (s *Sensibo) updateTargetTemperature() {
	state := s.Thermostat.TargetTemperature
	value := s.CurrentMeasurement.Temperature
	if s.TemperatureUnit == TemperatureUnitFahrenheit {
		value = CelciusToFahrenheit(value)
	}
	state.SetValue(float64(value))
}
