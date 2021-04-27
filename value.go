package config

import (
	"time"
)

type Value interface {
	ReadValue
	ParseValue
	UnmarshalValue
}

type UnmarshalValue interface {
	Unmarshal(val interface{}) error
}

type ReadValue interface {
	String() string
	Int() int
	Int64() int64
	Uint() uint
	Uint64() uint64
	Float64() float64
	Bool() bool
	Duration() time.Duration
	Time() time.Time
}

type ParseValue interface {
	ParseString() (string, error)
	ParseInt() (int, error)
	ParseInt64() (int64, error)
	ParseUint() (uint, error)
	ParseUint64() (uint64, error)
	ParseFloat64() (float64, error)
	ParseBool() (bool, error)
	ParseDuration() (time.Duration, error)
	ParseTime() (time.Time, error)
}
