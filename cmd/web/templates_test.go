package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	t.Parallel()
	tm := time.Date(2024, time.March, 17, 10, 15, 0, 0, time.UTC)
	hd := humanDate(tm)
	want := "17 Mar 2024 at 10:15"

	if hd != want {
		t.Errorf("got %q; want %q", hd, want)
	}
}
