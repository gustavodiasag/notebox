package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	td := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
	fd := fmtDate(td)

	if fd != "17 Mar, 2022" {
		t.Errorf("got %q, want %q", fd, "17 Mar, 2022")
	}
}
