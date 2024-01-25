package require

import "testing"

func Fail(t *testing.T, msg string, args ...interface{}) {
	t.Helper()
	t.Errorf(msg, args...)
	t.FailNow()
}
