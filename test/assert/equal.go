package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	t.Error(msgAndArgs...)

	return false
}

func Equalf(t *testing.T, expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	t.Errorf(msg, args...)

	return false
}
