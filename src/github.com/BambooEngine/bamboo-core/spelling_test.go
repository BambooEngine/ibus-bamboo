package bamboo

import (
	"testing"
)

func TestGenerateDumpSoundsFromWord(t *testing.T) {
	s := ParseDumpSoundsFromWord("chuyển")
	if len(s) != 6 {
		t.Errorf("Test length of chuyển, expected [6], got [%v]", len(s))
	}
}

func TestGenerateSoundsFromWord(t *testing.T) {
	s := ParseSoundsFromWord("chuyển")
	if len(s) != 6 {
		t.Errorf("Test length of chuyển, expected [6], got [%v]", len(s))
	}
	s2 := ParseSoundsFromWord("quyển")
	if len(s2) != 5 {
		t.Errorf("Test length of chuyển, expected [5], got [%v]", len(s2))
	}
	s3 := ParseSoundsFromWord("giá")
	if len(s3) != 3 {
		t.Errorf("Test length of chuyển, expected [3], got [%v]", len(s3))
	}
	s4 := ParseSoundsFromWord("gic1")
	if len(s4) != 4 {
		t.Errorf("Test length of chuyển, expected [4], got [%v]", len(s4))
	}
}
