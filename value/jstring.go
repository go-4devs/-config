package value

import (
	"time"

	"gitoa.ru/go-4devs/config"
)

var _ config.Value = (*JString)(nil)

type JString string

func (s JString) Unmarshal(v interface{}) error {
	return JUnmarshal(s.Bytes(), v)
}

func (s JString) ParseString() (string, error) {
	return s.String(), nil
}

func (s JString) ParseInt() (int, error) {
	return Atoi(s.String())
}

func (s JString) ParseInt64() (int64, error) {
	return ParseInt(s.String())
}

func (s JString) ParseUint() (uint, error) {
	u64, err := ParseUint(s.String())
	if err != nil {
		return 0, err
	}

	return uint(u64), nil
}

func (s JString) ParseUint64() (uint64, error) {
	return ParseUint(s.String())
}

func (s JString) ParseFloat64() (float64, error) {
	return ParseFloat(s.String())
}

func (s JString) ParseBool() (bool, error) {
	return ParseBool(s.String())
}

func (s JString) ParseDuration() (time.Duration, error) {
	return ParseDuration(s.String())
}

func (s JString) ParseTime() (time.Time, error) {
	return ParseTime(s.String())
}

func (s JString) Bytes() []byte {
	return []byte(s)
}

func (s JString) String() string {
	return string(s)
}

func (s JString) Int() int {
	in, _ := s.ParseInt()

	return in
}

func (s JString) Int64() int64 {
	in, _ := s.ParseInt64()

	return in
}

func (s JString) Uint() uint {
	in, _ := s.ParseUint()

	return in
}

func (s JString) Uint64() uint64 {
	in, _ := s.ParseUint64()

	return in
}

func (s JString) Float64() float64 {
	in, _ := s.ParseFloat64()

	return in
}

func (s JString) Bool() bool {
	in, _ := s.ParseBool()

	return in
}

func (s JString) Duration() time.Duration {
	in, _ := s.ParseDuration()

	return in
}

func (s JString) Time() time.Time {
	in, _ := s.ParseTime()

	return in
}
