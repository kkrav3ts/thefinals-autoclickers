package timing

import (
	"testing"
	"time"
)

func TestGenerateClickDelaysBounds(t *testing.T) {
	const n = 5000
	got, err := GenerateClickDelays(n, 60, 5, 50, 80)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != n {
		t.Fatalf("len = %d, want %d", len(got), n)
	}

	var sum float64
	for _, d := range got {
		ms := float64(d) / float64(time.Millisecond)
		if ms < 50 || ms > 80 {
			t.Fatalf("delay %vms outside [50, 80]", ms)
		}
		sum += ms
	}
	avg := sum / n
	if avg < 55 || avg > 65 {
		t.Fatalf("average %vms, want it near 60", avg)
	}
}

func TestGenerateClickDelaysRejectsBadInput(t *testing.T) {
	if _, err := GenerateClickDelays(0, 1, 1, 1, 2); err == nil {
		t.Fatal("expected error for non-positive count")
	}
	if _, err := GenerateClickDelays(1, 1, -0.1, 1, 2); err == nil {
		t.Fatal("expected error for negative stddev")
	}
	if _, err := GenerateClickDelays(1, 1, 1, 5, 2); err == nil {
		t.Fatal("expected error when min exceeds max")
	}
}
