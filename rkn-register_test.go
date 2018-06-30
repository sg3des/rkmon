package main

import (
	"log"
	"net"
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
	t.Log(testReg.TotalIP())
}

func BenchmarkRegister(b *testing.B) {
	ip := net.IP{255, 255, 255, 255}
	for i := 0; i < b.N; i++ {
		testReg.LookupIP(ip)
	}
}
