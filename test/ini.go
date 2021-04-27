package test

import (
	"log"

	"gopkg.in/ini.v1"
)

func NewINI() *ini.File {
	f, err := ini.Load(FixturePath("config.ini"))
	if err != nil {
		log.Fatal(err)
	}

	return f
}
