package inmem

import (
	"github.com/edkvm/ctrl/action"
	"sync"
	"time"
)

type actionRepo struct {
	mtx sync.RWMutex
	actions map[string]*action.Action
}

func NewActionRepo() *actionRepo {
	return &actionRepo{
		actions: make(map[string]*action.Action, 0),
	}
}

func (r *actionRepo) FindAll() []*action.Action {
	return nil
}

type triggerRepo struct {
	mtx sync.RWMutex
	triggers map[action.ScheduleID]*action.Schedule
}

func NewTriggerRepo() *triggerRepo {
	return &triggerRepo{
		triggers: make(map[action.ScheduleID]*action.Schedule, 0),
	}
}

func (r *triggerRepo) Store(t *action.Schedule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.triggers[t.ID] = t
	return nil
}

func (r *triggerRepo) FindClosestScheduleTrigger(dur time.Duration) *action.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var item *action.Schedule
	minDur := dur
	for _, val := range r.triggers {
		curDur := val.When.Sub(time.Now())
		if curDur < minDur {
			minDur = curDur
			item = val
		}
	}

	return item
}

func (r *triggerRepo) FindAllScheduleTriggers(time time.Time) []*action.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*action.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.When == time {
			items = append(items, val)
		}
	}

	return items
}

type statsRepo struct {
	mtx sync.RWMutex
	stats map[action.RegID][]*action.Stat
}


func NewStatsRepo() *statsRepo {
	return &statsRepo{
		stats: make(map[action.RegID][]*action.Stat, 0),
	}
}

func (r *statsRepo) Store(stat *action.Stat) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.stats[stat.ActionRegID]; !ok {
		r.stats[stat.ActionRegID] = make([]*action.Stat, 0)
	}
	r.stats[stat.ActionRegID] = append(r.stats[stat.ActionRegID], stat)

	return nil
}

func (r *statsRepo) FindLast(actionID action.RegID) (*action.Stat, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if list, ok := r.stats[actionID]; !ok {
		last := len(list) - 1
		return list[last], nil
	}

	return nil, action.ErrMissingStats
}


