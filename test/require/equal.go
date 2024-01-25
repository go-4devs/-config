package require

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Error(msgAndArgs...)
	t.FailNow()
}

func Equalf(t *testing.T, expected interface{}, actual interface{}, msg string, args ...interface{}) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf(msg, args...)
	t.FailNow()
}
