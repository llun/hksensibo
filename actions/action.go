package actions

import (
	"github.com/brutella/hc/log"
	"github.com/llun/sensibo-golang"

	"container/list"
	"time"
)

type Store interface {
	CurrentAcState() sensibo.AcState
	CurrentMeasurement() sensibo.Measurement
	UpdateAcState(state sensibo.AcState)
	UpdateMeasurement(measurement sensibo.Measurement)
}

type Action interface {
	Run(api *sensibo.Sensibo, pod sensibo.Pod, store Store)
	Name() string
	RemoveDuplicateCommand() bool
}

type Worker struct {
	actions  *list.List
	tickerCh <-chan time.Time

	api   *sensibo.Sensibo
	pod   sensibo.Pod
	store Store
}

func NewWorker(api *sensibo.Sensibo, pod sensibo.Pod, store Store) *Worker {
	return &Worker{
		actions:  list.New(),
		tickerCh: time.Tick(1 * time.Second),
		api:      api,
		pod:      pod,
		store:    store,
	}
}

func (w *Worker) Run() {
	for range w.tickerCh {
		if w.actions.Len() == 0 {
			continue
		}

		head := w.actions.Front()
		action, ok := w.actions.Remove(head).(Action)
		if !ok {
			continue
		}

		log.Debug.Printf("Sensibo run %v", action.Name())
		action.Run(w.api, w.pod, w.store)
	}
}

func (w *Worker) AddAction(action Action) {
	if w.actions.Len() > 0 {
		tail := w.actions.Back()
		lastAction, ok := tail.Value.(Action)
		if !ok {
			return
		}

		if lastAction.RemoveDuplicateCommand() &&
			lastAction.Name() == action.Name() {
			w.actions.Remove(tail)
		}
	}
	w.actions.PushBack(action)
}
