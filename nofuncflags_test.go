package main

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f := new("fixtures", os.Stdout)
	f.parse()
	f.print()
	if len(f.res) != 1 {
		t.Fail()
	}
}
