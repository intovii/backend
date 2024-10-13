package entities

import (
	"time"
)


type Workout struct {
	StartTime string
	EndTime   string
	FormatID  uint64
	UsersID   []uint64
	TrainerID uint64
	FilialID  uint64
	Status    string
	Date      string
}

type Day struct {
	Workouts []Workout
	Date     string
}

type SchedulerGetter struct {
	ID    uint64
	Start time.Time
	End   time.Time
}
