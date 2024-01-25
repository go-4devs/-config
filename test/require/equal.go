package require

import (
	"testing"

	"gitoa.ru/go-4devs/config/test/assert"
)

func Equal(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if assert.Equal(t, expected, actual, msgAndArgs...) {
		return
	}

	t.FailNow()
}

func Equalf(t *testing.T, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	t.Helper()

	if assert.Equalf(t, expected, actual, msg, args...) {
		return
	}

	t.Errorf(msg, args...)
	t.FailNow()
}
