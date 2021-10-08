package utils

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

func Time2MilliTimestamp(input null.Time) TimeStampInt64 {
	if !input.Valid {
		return NewTimeStampInt64(input.Time.UnixNano()/int64(time.Millisecond), false)
	} else {
		return NewTimeStampInt64(input.Time.UnixNano()/int64(time.Millisecond), true)
	}

}

func MilliTimestamp2Time(input int64) time.Time {
	return time.Unix(0, input*int64(time.Millisecond))
}
