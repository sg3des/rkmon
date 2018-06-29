package main

import (
	"net"
	"os"
	"testing"
)

func TestSendAlertTemplate(t *testing.T) {
	s := SMTP{
		Host: "smtp.gmail.com:587",
		User: os.Getenv("TEST_USER"),
		Pass: os.Getenv("TEST_PASS"),
	}

	if s.User == "" || s.Pass == "" {
		t.Log("set environment variables user and password, export TEST_USER=email@gmail.com; export TEST_PASS=password")
		t.SkipNow()
	}

	err := s.SendAlertTemplate(os.Getenv("TEST_USER"), []net.IP{net.IP{127, 0, 0, 1}, net.IP{255, 255, 255, 255}})
	if err != nil {
		t.Error(err)
	}
}
