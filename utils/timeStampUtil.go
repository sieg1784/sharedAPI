package utils

import (
	"database/sql"
	"encoding/json"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type TimeStampInt64 struct {
	sql.NullInt64
}

func (v TimeStampInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *TimeStampInt64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int64 = *s
	} else {
		v.Valid = false
	}
	return nil
}

// NewTime creates a new Time.
func NewTimeStampInt64(t int64, valid bool) TimeStampInt64 {
	return TimeStampInt64{
		NullInt64: sql.NullInt64{
			Int64: t,
			Valid: valid,
		},
	}
}
