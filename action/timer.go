package action

import (
	"log"
	"sync"
	"time"
)

type ActionTimer struct {
	mtx sync.RWMutex
	scheduled map[ScheduleID]*time.Timer
}

func NewActionTimer() *ActionTimer {
	return &ActionTimer{
		scheduled: make(map[ScheduleID]*time.Timer, 0),
	}
}

func (at *ActionTimer) RemoveSchedule(id ScheduleID) {
	at.mtx.Lock()
	defer at.mtx.Unlock()

	if val, ok := at.scheduled[id]; ok {
		val.Stop()
	}

	delete(at.scheduled, id)
}

func (at *ActionTimer) AddSchedule (sched *Schedule, fn func(name string, params map[string]interface{}) (interface{}, error)) error {
	log.Println("add sched")
	t := time.AfterFunc(time.Duration(sched.Interval) * time.Second, func() {
		if sched.Recurring {
			at.AddSchedule(sched, fn)
		}

		if _, ok := at.scheduled[sched.ID]; !ok {
			return
		}
		_, err := fn(sched.Action, sched.Params)
		if err != nil {

		}
	})

	at.mtx.Lock()
	defer at.mtx.Unlock()

	at.scheduled[sched.ID] = t

	return nil
}
