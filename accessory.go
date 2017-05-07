package hksensibo

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/service"
	"github.com/llun/hksensibo/actions"
	"github.com/llun/sensibo-golang"

	"time"
)

type Sensibo struct {
	*accessory.Accessory

	Thermostat  *service.Thermostat
	state       sensibo.AcState
	measurement sensibo.Measurement

	worker    *actions.Worker
	pollingCh <-chan time.Time
}

func (s *Sensibo) PollingState() {
	s.worker.AddAction(actions.NewGetAcState())
	s.worker.AddAction(actions.NewGetMeasurement())
	for range s.pollingCh {
		s.worker.AddAction(actions.NewGetAcState())
		s.worker.AddAction(actions.NewGetMeasurement())
	}
}

func (s *Sensibo) UpdateAcState(state sensibo.AcState) {
	s.state = state

	targetState := characteristic.TargetHeatingCoolingStateOff
	currentState := characteristic.CurrentHeatingCoolingStateOff
	if state.On {
		switch state.Mode {
		case DRY_MODE:
			targetState = characteristic.TargetHeatingCoolingStateHeat
			currentState = characteristic.CurrentHeatingCoolingStateHeat
		case COOL_MODE:
			targetState = characteristic.TargetHeatingCoolingStateCool
			currentState = characteristic.TargetHeatingCoolingStateCool
		default:
			targetState = characteristic.TargetHeatingCoolingStateAuto
			currentState = characteristic.TargetHeatingCoolingStateCool
		}
	}

	s.Thermostat.CurrentHeatingCoolingState.UpdateValue(currentState)
	s.Thermostat.TargetHeatingCoolingState.UpdateValue(targetState)
	s.Thermostat.TargetTemperature.UpdateValue(float64(state.TargetTemperature))
}

func (s *Sensibo) UpdateMeasurement(measurement sensibo.Measurement) {
	s.measurement = measurement
	s.Thermostat.CurrentTemperature.UpdateValue(measurement.Temperature)
}

func (s *Sensibo) CurrentAcState() sensibo.AcState {
	return s.state
}

func (s *Sensibo) CurrentMeasurement() sensibo.Measurement {
	return s.measurement
}

func NewSensibo(pod sensibo.Pod, api *sensibo.Sensibo) *Sensibo {
	info := accessory.Info{
		Name:         "Sensibo",
		Manufacturer: "Sensibo",
		SerialNumber: pod.ID,
		Model:        pod.Room.Name,
	}

	acc := Sensibo{
		Thermostat: service.NewThermostat(),
		pollingCh:  time.Tick(60 * time.Second),
	}
	acc.Accessory = accessory.New(info, accessory.TypeThermostat)
	acc.AddService(acc.Thermostat.Service)
	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(acc.onTargetTemperatureUpdate)
	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(acc.onHeatingCoolingStateUpdate)

	worker := actions.NewWorker(api, pod, &acc)
	acc.worker = worker
	go worker.Run()
	go acc.PollingState()
	return &acc
}

func Lookup(key string) []*Sensibo {
	api := sensibo.NewSensibo(key)
	pods, err := api.GetPods()
	if err != nil {
		log.Info.Fatal(err)
		return nil
	}

	var services []*Sensibo = make([]*Sensibo, len(pods))
	for index, pod := range pods {
		services[index] = NewSensibo(pod, api)
	}

	return services
}
