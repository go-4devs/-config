package value

import (
	"time"

	"gitoa.ru/go-4devs/config"
)

var _ config.Value = (*JBytes)(nil)

type JBytes []byte

func (s JBytes) Unmarshal(v interface{}) error {
	return JUnmarshal(s.Bytes(), v)
}

func (s JBytes) ParseString() (string, error) {
	return s.String(), nil
}

func (s JBytes) ParseInt() (int, error) {
	return Atoi(s.String())
}

func (s JBytes) ParseInt64() (int64, error) {
	return ParseInt(s.String())
}

func (s JBytes) ParseUint() (uint, error) {
	u64, err := ParseUint(s.String())
	if err != nil {
		return 0, err
	}

	return uint(u64), nil
}

func (s JBytes) ParseUint64() (uint64, error) {
	return ParseUint(s.String())
}

func (s JBytes) ParseFloat64() (float64, error) {
	return ParseFloat(s.String())
}

func (s JBytes) ParseBool() (bool, error) {
	return ParseBool(s.String())
}

func (s JBytes) ParseDuration() (time.Duration, error) {
	return ParseDuration(s.String())
}

func (s JBytes) ParseTime() (time.Time, error) {
	return ParseTime(s.String())
}

func (s JBytes) Bytes() []byte {
	return []byte(s)
}

func (s JBytes) String() string {
	return string(s)
}

func (s JBytes) Int() int {
	in, _ := s.ParseInt()

	return in
}

func (s JBytes) Int64() int64 {
	in, _ := s.ParseInt64()

	return in
}

func (s JBytes) Uint() uint {
	in, _ := s.ParseUint()

	return in
}

func (s JBytes) Uint64() uint64 {
	in, _ := s.ParseUint64()

	return in
}

func (s JBytes) Float64() float64 {
	in, _ := s.ParseFloat64()

	return in
}

func (s JBytes) Bool() bool {
	in, _ := s.ParseBool()

	return in
}

func (s JBytes) Duration() time.Duration {
	in, _ := s.ParseDuration()

	return in
}

func (s JBytes) Time() time.Time {
	in, _ := s.ParseTime()

	return in
}

func (s JBytes) IsEquals(in config.Value) bool {
	return s.String() == in.String()
}
