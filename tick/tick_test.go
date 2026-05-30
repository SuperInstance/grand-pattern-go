package tick

import (
	"math"
	"testing"
)

func TestTickSchedule(t *testing.T) {
	ts := NewTickSchedule()
	called := 0
	ts.At(5, func() { called++ })
	ts.Fire(5)
	if called != 1 {
		t.Fatalf("called = %d, want 1", called)
	}
}

func TestTempo(t *testing.T) {
	sec := Tempo(120.0)
	if math.Abs(sec-0.5) > 1e-9 {
		t.Fatalf("tempo = %f, want 0.5", sec)
	}
}
