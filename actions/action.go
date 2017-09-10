package actions

import (
	"github.com/llun/sensibo-golang"
)

const RETRY_COUNT = 3

type Store interface {
	CurrentAcState() sensibo.AcState
	CurrentMeasurement() sensibo.Measurement
	UpdateAcState(state sensibo.AcState)
	UpdateMeasurement(measurement sensibo.Measurement)
}
