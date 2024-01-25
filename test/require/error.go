package require

import (
	"testing"
)

func NoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()

	if err != nil {
		t.Error(msgAndArgs...)
		t.FailNow()
	}
}

func NoErrorf(t *testing.T, err error, msg string, args ...interface{}) {
	t.Helper()

	if err != nil {
		t.Errorf(msg, args...)
		t.FailNow()
	}
}
