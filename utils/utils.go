package utils

import (
	"io/ioutil"
	"log"
)

func Check(err error, info string) {
	if err != nil {
		log.Fatalf("%s: %q", info, err)
	}
}

func ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	Check(err, "Error reading file")
	return string(data)
}
