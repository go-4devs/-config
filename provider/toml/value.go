package toml

import (
	"encoding/json"
	"fmt"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
)

type Value struct {
	value.Value
}

func (s Value) Int() int {
	v, _ := s.ParseInt()

	return v
}

func (s Value) ParseInt() (int, error) {
	v, err := s.ParseInt64()
	if err != nil {
		return 0, fmt.Errorf("toml failed parce int: %w", err)
	}

	return int(v), nil
}

func (s Value) Unmarshal(target interface{}) error {
	b, err := json.Marshal(s.Raw())
	if err != nil {
		return fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return nil
}
