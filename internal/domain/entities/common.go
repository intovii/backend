package entities

import (
	"database/sql"
	"time"
)


func NewNullString(s string) *sql.NullString {
	if len(s) == 0 {
		return &sql.NullString{}
	} else {
		return &sql.NullString{
			String: s,
			Valid: true,
		}
	}
}

func NewNullTime(t *time.Time) *sql.NullTime {
	if t == nil {
		return &sql.NullTime{}
	} else {
		return &sql.NullTime{
			Time: *t,
			Valid: true,
		}
	}
}
