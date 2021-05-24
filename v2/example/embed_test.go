package example

import (
	"log"
	"testing"
)

func TestDemoSpecifyFile(t *testing.T) {
	if err := demoSpecifyFile(); err != nil {
		log.Fatal(err)
	}
}

func TestDemoDir(t *testing.T) {
	// I think this is a general case more than any others, that you just prepare a directory and put all language on it.
	demoDir()
}
