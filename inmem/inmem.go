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

type scheduleRepo struct {
	mtx sync.RWMutex
	triggers map[action.ScheduleID]*action.Schedule
}

func NewTriggerRepo() *scheduleRepo {
	return &scheduleRepo{
		triggers: make(map[action.ScheduleID]*action.Schedule, 0),
	}
}

func (r *scheduleRepo) Store(t *action.Schedule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.triggers[t.ID] = t
	return nil
}

func (r *scheduleRepo) FindNext(dur time.Duration) *action.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var item *action.Schedule
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

func (r *scheduleRepo) FindAllByTime(time time.Time) []*action.Schedule {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	items := make([]*action.Schedule, 0, len(r.triggers))
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

func (r *statsRepo) FindByID(actionID string) ([]*action.Stat, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	list := make([]*action.Stat, 0)

	for _, val := range r.stats {
		if val.ActionName == actionID {
			list = append(list, val)
		}
	}

	return list, action.ErrMissingStats
}


