package action

import (
	"time"
	ctrlID "github.com/edkvm/ctrl/pkg/id"
)

type Stat struct {
	ID     string
	Action string
	Start  time.Time
	End    time.Time
	Status bool
}

func NewStat(name string, start, end time.Time, status bool) *Stat{
	return &Stat{
		ID: ctrlID.GenULID(),
		Action: name,
		Start:  start,
		End:    time.Now(),
		Status: status,
	}
}

type StatsRepo interface {
	Store(stat *Stat) error
	FindAll() []*Stat
	FindByAction(actionID string) []*Stat
}
