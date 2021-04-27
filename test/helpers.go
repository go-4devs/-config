package test

import (
	"path/filepath"
	"runtime"
)

func FixturePath(file string) string {
	path := "fixture/"

	_, filename, _, ok := runtime.Caller(0)
	if ok {
		path = filepath.Dir(filename) + "/" + path
	}

	return path + file
}
