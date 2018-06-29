package main

import (
	"testing"
)

var (
	testrep *Repository
)

func TestNewRepository(t *testing.T) {
	rep, err := NewRepository("https://github.com/sg3des/stob")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testrep = rep
}

func TestGet(t *testing.T) {
	if testrep == nil {
		t.SkipNow()
	}

	time, err := testrep.Get()
	if err != nil {
		t.Error(err)
	}

	if time.IsZero() {
		t.Error("commit time is nil")
	}
}

func TestOpenFile(t *testing.T) {
	if testrep == nil {
		t.SkipNow()
	}

	f, err := testrep.OpenFile("stob.go")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	f.Close()
}
