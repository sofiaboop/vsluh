package domain

import "testing"

func TestModeDurations(t *testing.T) {
	cases := []struct {
		mode      Mode
		prepSec   int
		speakSec  int
		wantValid bool
	}{
		{ModeSpontaneous, 30, 60, true},
		{ModeDeep, 120, 180, true},
		{Mode("unknown"), 30, 60, false},
	}

	for _, c := range cases {
		prep, speak := c.mode.Durations()
		if prep != c.prepSec || speak != c.speakSec {
			t.Errorf("%s: got (%d, %d), want (%d, %d)", c.mode, prep, speak, c.prepSec, c.speakSec)
		}
		if got := c.mode.Valid(); got != c.wantValid {
			t.Errorf("%s: Valid() = %v, want %v", c.mode, got, c.wantValid)
		}
	}
}
