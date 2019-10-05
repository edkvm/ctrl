package inmem

import (
	"github.com/edkvm/ctrl/trigger"
	"sync"
	"time"
)

type triggerRepo struct {
	mtx sync.RWMutex
	triggers map[trigger.ScheduleID]*trigger.Schedule
}


func (r *triggerRepo) Store(t *trigger.Schedule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.triggers[t.ID] = t
	return nil
}

func (r *triggerRepo) FindClosestScheduleTrigger(dur time.Duration) *trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var item *trigger.Schedule
	minDur := dur
	for _, val := range r.triggers {
		curDur := val.Time.Sub(time.Now())
		if curDur < minDur {
			minDur = curDur
			item = val
		}
	}

	return item
}

func (r *triggerRepo) FindAllScheduleTriggers(time time.Time) []*trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*trigger.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.Time == time {
			items = append(items, val)
		}
	}

	return items
}


