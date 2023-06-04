package value

import (
	"time"
)

type Empty struct {
	Err error
}

func (e Empty) Unmarshal(_ interface{}) error {
	return e.Err
}

func (e Empty) ParseString() (string, error) {
	return "", e.Err
}

func (e Empty) ParseInt() (int, error) {
	return 0, e.Err
}

func (e Empty) ParseInt64() (int64, error) {
	return 0, e.Err
}

func (e Empty) ParseUint() (uint, error) {
	return 0, e.Err
}

func (e Empty) ParseUint64() (uint64, error) {
	return 0, e.Err
}

func (e Empty) ParseFloat64() (float64, error) {
	return 0, e.Err
}

func (e Empty) ParseBool() (bool, error) {
	return false, e.Err
}

func (e Empty) ParseDuration() (time.Duration, error) {
	return 0, e.Err
}

func (e Empty) ParseTime() (time.Time, error) {
	return time.Time{}, e.Err
}

func (e Empty) String() string {
	return ""
}

func (e Empty) Int() int {
	return 0
}

func (e Empty) Int64() int64 {
	return 0
}

func (e Empty) Uint() uint {
	return 0
}

func (e Empty) Uint64() uint64 {
	return 0
}

func (e Empty) Float64() float64 {
	return 0
}

func (e Empty) Bool() bool {
	return false
}

func (e Empty) Duration() time.Duration {
	return 0
}

func (e Empty) Time() time.Time {
	return time.Time{}
}
