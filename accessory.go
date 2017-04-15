package sensibo

import (
  "github.com/brutella/hc/accessory"
  "github.com/brutella/hc/log"
  "github.com/brutella/hc/service"
  "github.com/llun/sensibo-golang"

  "time"
)

type TemperatureUnit string

const (
  TemperatureUnitCelcius    TemperatureUnit = "C"
  TemperatureUnitFahrenheit                 = "F"
)

type Sensibo struct {
  *accessory.Accessory

  Thermostat         *service.Thermostat
  CurrentState       sensibo.AcState
  CurrentMeasurement sensibo.Measurement
  TemperatureUnit    TemperatureUnit

  api *sensibo.Sensibo
  pod sensibo.Pod
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

func NewSensibo(pod sensibo.Pod, api *sensibo.Sensibo) *Sensibo {
  info := accessory.Info{
    Name:         "Sensibo",
    Manufacturer: "Sensibo",
    SerialNumber: pod.ID,
    Model:        pod.Room.Name,
  }

  acc := Sensibo{
    TemperatureUnit: TemperatureUnitCelcius,
    pod:             pod,
    api:             api,
  }
  acc.Accessory = accessory.New(info, accessory.TypeThermostat)
  acc.Thermostat = service.NewThermostat()
  acc.AddService(acc.Thermostat.Service)

  states, err := api.GetAcStates(pod.ID)
  if err != nil {
    log.Info.Fatal(err)
  }
  if len(states) != 0 {
    acc.CurrentState = states[0].AcState
  }

  measurements, err := api.GetMeasurements(pod.ID)
  if err != nil {
    log.Info.Fatal(err)
  }
  if len(measurements) != 0 {
    acc.CurrentMeasurement = measurements[0]
  }

  acc.setup()
  return &acc
}

func (s *Sensibo) setup() {
  s.setupTemperatureDisplayUnits()
  s.setupTargetHeatingCoolingState()
  s.setupTargetTemperature()
  s.updateCurrentTemperature()
  s.updateCurrentHeatingCoolingState()

  go func() {
    c := time.Tick(30 * time.Second)
    for range c {
      s.updateHomeKitFromState()
    }
  }()
}

func (s *Sensibo) updateHomeKitFromState() {
  s.updateTemperatureDisplayUnits()
  s.updateCurrentTemperature()
  s.updateTargetTemperature()
  s.updateCurrentHeatingCoolingState()
  s.updateTargetHeatingCoolingState()
}
