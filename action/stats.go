package action

import "time"

type Stat struct {
	ID         string
	ActionName string
	Start      time.Time
	End        time.Time
	Status     bool
}

type StatsRepo interface {
	Store(stat *Stat) error
	FindAll() []*Stat
	FindByID(actionID string) ([]*Stat, error)
}
