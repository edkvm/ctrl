package inmem

import (
	"fmt"
	"github.com/edkvm/ctrl/action"
	"github.com/edkvm/ctrl/trigger"
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

type scheduleRepo struct {
	mtx sync.RWMutex
	triggers map[trigger.ScheduleID]*trigger.Schedule
}

func NewScheduleRepo() *scheduleRepo {
	return &scheduleRepo{
		triggers: make(map[trigger.ScheduleID]*trigger.Schedule, 0),
	}
}

func (r *scheduleRepo) Store(t *trigger.Schedule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.triggers[t.ID] = t
	return nil
}

func (r *scheduleRepo) FindNext(dur time.Duration) *trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var item *trigger.Schedule
	minDur := dur
	for _, val := range r.triggers {
		curDur := val.Start.Sub(time.Now())
		if curDur < minDur {
			minDur = curDur
			item = val
		}
	}

	return item
}

func (r *scheduleRepo) Find(id trigger.ScheduleID) (*trigger.Schedule, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.triggers[id]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("not found")
}

func (r *scheduleRepo) FindAllByAction(name string) []*trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*trigger.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.Action == name {
			items = append(items, val)
		}
	}

	return items
}

func (r *scheduleRepo) FindAllByTime(time time.Time) []*trigger.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*trigger.Schedule, 0, len(r.triggers))
	for _, val := range r.triggers {
		if val.Start == time {
			items = append(items, val)
		}
	}

	return items
}


type statsRepo struct {
	mtx sync.RWMutex
	stats map[string]*action.Stat
}

func NewStatsRepo() *statsRepo {
	return &statsRepo{
		stats: make(map[string]*action.Stat, 0),
	}
}

func (r *statsRepo) Store(stat *action.Stat) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.stats[stat.ID] = stat

	return nil
}

func (r *statsRepo) FindAll() []*action.Stat {

	stats := make([]*action.Stat, 0, len(r.stats))

	for _, v := range r.stats {
		stats = append(stats, v)
	}

	return stats
}

func (r *statsRepo) FindByAction(actionID string) []*action.Stat {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	list := make([]*action.Stat, 0)

	for _, val := range r.stats {
		if val.Action == actionID {
			list = append(list, val)
		}
	}

	return list
}


