package value

import (
	"encoding/json"
	"fmt"
	"time"

	"gitoa.ru/go-4devs/config"
)

var _ config.Value = (*Value)(nil)

type Value struct {
	Val interface{}
}

func (s Value) Int() int {
	v, _ := s.ParseInt()

	return v
}

func (s Value) Int64() int64 {
	v, _ := s.ParseInt64()

	return v
}

func (s Value) Uint() uint {
	v, _ := s.ParseUint()

	return v
}

func (s Value) Uint64() uint64 {
	v, _ := s.ParseUint64()

	return v
}

func (s Value) Float64() float64 {
	in, _ := s.ParseFloat64()

	return in
}

func (s Value) String() string {
	v, _ := s.ParseString()

	return v
}

func (s Value) Bool() bool {
	v, _ := s.ParseBool()

	return v
}

func (s Value) Duration() time.Duration {
	v, _ := s.ParseDuration()

	return v
}

func (s Value) Raw() interface{} {
	return s.Val
}

func (s Value) Time() time.Time {
	v, _ := s.ParseTime()

	return v
}

func (s Value) Unmarshal(target interface{}) error {
	if v, ok := s.Raw().([]byte); ok {
		err := json.Unmarshal(v, target)
		if err != nil {
			return fmt.Errorf("%w: %s", config.ErrInvalidValue, err)
		}

		return nil
	}

	return config.ErrInvalidValue
}

func (s Value) ParseInt() (int, error) {
	if r, ok := s.Raw().(int); ok {
		return r, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseInt64() (int64, error) {
	if r, ok := s.Raw().(int64); ok {
		return r, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseUint() (uint, error) {
	if r, ok := s.Raw().(uint); ok {
		return r, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseUint64() (uint64, error) {
	if r, ok := s.Raw().(uint64); ok {
		return r, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseFloat64() (float64, error) {
	if r, ok := s.Raw().(float64); ok {
		return r, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseString() (string, error) {
	if r, ok := s.Raw().(string); ok {
		return r, nil
	}

	return "", config.ErrInvalidValue
}

func (s Value) ParseBool() (bool, error) {
	if b, ok := s.Raw().(bool); ok {
		return b, nil
	}

	return false, config.ErrInvalidValue
}

func (s Value) ParseDuration() (time.Duration, error) {
	if b, ok := s.Raw().(time.Duration); ok {
		return b, nil
	}

	return 0, config.ErrInvalidValue
}

func (s Value) ParseTime() (time.Time, error) {
	if b, ok := s.Raw().(time.Time); ok {
		return b, nil
	}

	return time.Time{}, config.ErrInvalidValue
}
