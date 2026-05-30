package jepa

import (
	"math"
	"testing"
)

func TestPerceive(t *testing.T) {
	j := &Jepa{Window: 3}
	j.Perceive([]float64{1, 2, 3}, 1)
	j.Perceive([]float64{2, 3, 4}, 2)
	if len(j.DbIn.Entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(j.DbIn.Entries))
	}
}

func TestPredict(t *testing.T) {
	j := &Jepa{Window: 3}
	j.Perceive([]float64{2, 2}, 1)
	j.Perceive([]float64{4, 4}, 2)
	pred := j.Predict()
	if len(pred) != 2 {
		t.Fatalf("pred dims = %d, want 2", len(pred))
	}
	// Average of (2,2) and (4,4) → (3,3)
	if math.Abs(pred[0]-3.0) > 1e-9 || math.Abs(pred[1]-3.0) > 1e-9 {
		t.Fatalf("pred = %v, want [3 3]", pred)
	}
}

func TestSurprise(t *testing.T) {
	s := Surprise([]float64{0, 0}, []float64{3, 4})
	if math.Abs(s-5) > 1e-9 {
		t.Fatalf("surprise = %f, want 5", s)
	}
}

func TestConservation(t *testing.T) {
	j := &Jepa{Window: 3}
	j.Perceive([]float64{1}, 1)
	j.Perceive([]float64{2}, 2)
	j.Predict()
	j.RecordPrediction([]float64{1.5}, []float64{2})
	if !j.CheckConservation(2) {
		t.Fatal("conservation check failed unexpectedly")
	}
}

func TestGC(t *testing.T) {
	j := &Jepa{Window: 3}
	j.RecordPrediction([]float64{0}, []float64{0})   // surprise 0
	j.RecordPrediction([]float64{0}, []float64{10})  // surprise 10
	j.GC(1.0)
	if len(j.DbOut.Entries) != 1 {
		t.Fatalf("after GC: %d entries, want 1", len(j.DbOut.Entries))
	}
}
