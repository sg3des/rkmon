package main

import (
	"log"
	"os"
	"testing"
)

var testReg *Register

func init() {
	testReg = NewRegister()
	log.SetFlags(log.Lshortfile)
}

func TestRegisterLoad(t *testing.T) {
	f, err := os.Open("z-i/dump.csv")
	if err != nil {
		t.Fatal("failed open dump.csv file,", err)
	}

	testReg.Load(f)

	if testReg.TotalIP() == 0 {
		t.Error("failed parse dump.csv")
	}
}
