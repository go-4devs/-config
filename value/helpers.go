package value

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gitoa.ru/go-4devs/config"
)

func ParseDuration(raw string) (time.Duration, error) {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return d, nil
}

func ParseInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return i, nil
}

func ParseUint(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return i, nil
}

func Atoi(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return i, nil
}

func ParseTime(s string) (time.Time, error) {
	i, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return i, nil
}

func ParseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return f, nil
}

func ParseBool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return b, nil
}

func JUnmarshal(b []byte, v interface{}) error {
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return nil
}
