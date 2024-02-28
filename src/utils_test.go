package main

import (
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
