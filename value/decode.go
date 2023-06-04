// nolint: nonamedreturns
package value

import (
	"time"

	"gitoa.ru/go-4devs/config"
)

var _ config.Value = (*Decode)(nil)

type Decode func(v interface{}) error

func (s Decode) Unmarshal(v interface{}) error {
	return s(v)
}

func (s Decode) ParseInt() (v int, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseInt64() (v int64, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseUint() (v uint, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseUint64() (v uint64, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseFloat64() (v float64, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseString() (v string, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseDuration() (v time.Duration, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) ParseTime() (v time.Time, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) Int() int {
	in, _ := s.ParseInt()

	return in
}

func (s Decode) Int64() int64 {
	in, _ := s.ParseInt64()

	return in
}

func (s Decode) Uint() uint {
	in, _ := s.ParseUint()

	return in
}

func (s Decode) Uint64() uint64 {
	in, _ := s.ParseUint64()

	return in
}

func (s Decode) Float64() float64 {
	in, _ := s.ParseFloat64()

	return in
}

func (s Decode) Bytes() []byte {
	return []byte(s.String())
}

func (s Decode) String() string {
	v, _ := s.ParseString()

	return v
}

func (s Decode) Bool() bool {
	in, _ := s.ParseBool()

	return in
}

func (s Decode) ParseBool() (v bool, err error) {
	return v, s.Unmarshal(&v)
}

func (s Decode) Duration() time.Duration {
	in, _ := s.ParseDuration()

	return in
}

func (s Decode) Time() time.Time {
	in, _ := s.ParseTime()

	return in
}
