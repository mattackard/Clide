package clide

import (
	"testing"
	"time"
)

func TestGetKeyDelay(t *testing.T) {
	typeSpeed := 100
	humanize := 0.0
	if getKeyDelay(typeSpeed, humanize) != time.Duration(100)*time.Millisecond {
		t.Error("Key delay should be constant when humanize is 0")
	}

	humanize = 0.5
	random := getKeyDelay(typeSpeed, humanize)
	if random > time.Duration(150)*time.Millisecond || random < time.Duration(50)*time.Millisecond {
		t.Errorf("Keydelay expected to be between 50 and 150. Actual: %v", random)
	}

	humanize = 50.0
	//shouldn't error when humanize is greater than expected
	getKeyDelay(typeSpeed, humanize)

	humanize = -3.0
	//shouldn't error when humanize is negative
	getKeyDelay(typeSpeed, humanize)
}
