package hksensibo

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/log"
	"github.com/brutella/hc/service"

	"github.com/llun/hksensibo/actions"
	"github.com/llun/sensibo-golang"

	ba "github.com/llun/hkbridge/accessories"

	"net"
	"time"
)

type Sensibo struct {
	*accessory.Accessory

	Thermostat     *service.Thermostat
	HumiditySensor *service.HumiditySensor
	state          sensibo.AcState
	measurement    sensibo.Measurement

	worker    *ba.Worker
	pollingCh <-chan time.Time
	pod       sensibo.Pod
	api       *sensibo.Sensibo
}

func (s *Sensibo) PollingState() {
	s.worker.AddAction(actions.NewGetAcState(s.api, s.pod, s))
	s.worker.AddAction(actions.NewGetMeasurement(s.api, s.pod, s))
	for range s.pollingCh {
		s.worker.AddAction(actions.NewGetAcState(s.api, s.pod, s))
		s.worker.AddAction(actions.NewGetMeasurement(s.api, s.pod, s))
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
	s.HumiditySensor.CurrentRelativeHumidity.UpdateValue(measurement.Humidity)
}

func (s *Sensibo) CurrentAcState() sensibo.AcState {
	return s.state
}

func (s *Sensibo) CurrentMeasurement() sensibo.Measurement {
	return s.measurement
}

func NewSensibo(pod sensibo.Pod, api *sensibo.Sensibo, worker *ba.Worker) *Sensibo {
	info := accessory.Info{
		Name:         "Sensibo",
		Manufacturer: "Sensibo",
		SerialNumber: pod.ID,
		Model:        pod.Room.Name,
	}

	acc := Sensibo{
		Thermostat:     service.NewThermostat(),
		HumiditySensor: service.NewHumiditySensor(),
		pollingCh:      time.Tick(60 * time.Second),
		worker:         worker,
		api:            api,
		pod:            pod,
	}
	acc.Accessory = accessory.New(info, accessory.TypeThermostat)
	acc.AddService(acc.Thermostat.Service)
	acc.AddService(acc.HumiditySensor.Service)
	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(acc.onTargetTemperatureUpdate)
	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(acc.onHeatingCoolingStateUpdate)

	go acc.PollingState()
	return &acc
}

func Lookup(key string, worker *ba.Worker) []*Sensibo {
	api := sensibo.NewSensibo(key)
	pods, err := api.GetPods()
	if err != nil {
		log.Info.Fatal(err)
		return nil
	}

	var services []*Sensibo = make([]*Sensibo, len(pods))
	for index, pod := range pods {
		services[index] = NewSensibo(pod, api, worker)
	}

	return services
}

func AllAccessories(config ba.AccessoryConfig, iface *net.Interface, worker *ba.Worker) []*accessory.Accessory {
	option := config.Option

	key, ok := option["key"].(string)
	if !ok {
		log.Info.Println("Cannot read sensibo key")
		return nil
	}

	sensibos := Lookup(key, worker)
	sensiboAccessories := make([]*accessory.Accessory, len(sensibos))
	for idx, sensibo := range sensibos {
		sensiboAccessories[idx] = sensibo.Accessory
	}
	return sensiboAccessories
}
