package main

import (
	"os"
	"sort"
	"testing"
)

func TestSortStringList(t *testing.T) {
	var data = []string{"a", "ab", "ca"}
	sort.Sort(byString(data))
	if data[0] != "a" {
		t.Errorf("Sorting strings, expected %s, got %s", "a", data[0])
	}
	if data[1] != "ab" {
		t.Errorf("Sorting strings, expected %s, got %s", "ab", data[0])
	}
	if data[2] != "ca" {
		t.Errorf("Sorting strings, expected %s, got %s", "ca", data[0])
	}
}

func TestX11GetFocusWindowClassWithoutDisplay(t *testing.T) {
	display, hadDisplay := os.LookupEnv("DISPLAY")
	if err := os.Unsetenv("DISPLAY"); err != nil {
		t.Fatalf("failed to unset DISPLAY: %v", err)
	}
	defer func() {
		if hadDisplay {
			_ = os.Setenv("DISPLAY", display)
		} else {
			_ = os.Unsetenv("DISPLAY")
		}
	}()

	if wmClass := x11GetFocusWindowClass(); wmClass != "" {
		t.Fatalf("expected no window class without an X display, got %q", wmClass)
	}
}
