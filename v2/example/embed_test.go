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
