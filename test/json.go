package test

import (
	"io/ioutil"
	"log"
)

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(FixturePath(file))
	if err != nil {
		log.Fatal(err)
	}

	return data
}
