package wl

import (
	"math"
	"testing"
)

func TestFixedToFloat64(t *testing.T) {
	var f int32
	var d float64

	f = 0x012030
	d = fixedToFloat64(f)
	if d != 288.1875 {
		t.Fail()
	}

	f = -0x012030
	d = fixedToFloat64(f)
	if d != -288.1875 {
		t.Fail()
	}
}

func TestReverse(t *testing.T) {
	var d float64
	var f int32

	d = 3.1415
	f = float64ToFixed(d)
	dd := d - fixedToFloat64(f)
	if math.Abs(dd) > 0.001 {
		t.Fail()
	}
}
