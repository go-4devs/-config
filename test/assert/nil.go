package assert

import "testing"

func Nil(t *testing.T, data any, msgAndArgs ...interface{}) bool {
	t.Helper()

	if data != nil {
		t.Error(msgAndArgs...)

		return false
	}

	return true
}
