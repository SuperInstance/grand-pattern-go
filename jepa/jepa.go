package jepa

import (
	"math"
)

// Entry holds a vector observation with metadata.
type Entry struct {
	Data      []float64
	Timestamp uint64
	Surprise  float64
}

// VectorDb is a simple append-only vector store.
type VectorDb struct {
	Entries []Entry
}

// Jepa implements Joint-Embedding Predictive Architecture concepts:
// perceive input, predict next state, measure surprise.
type Jepa struct {
	DbIn   VectorDb
	DbOut  VectorDb
	Window int
	last   []float64
}

// Perceive records an observation into DbIn.
func (j *Jepa) Perceive(data []float64, ts uint64) {
	cp := make([]float64, len(data))
	copy(cp, data)
	j.DbIn.Entries = append(j.DbIn.Entries, Entry{Data: cp, Timestamp: ts})
}

// Predict returns a naive moving-average prediction based on the last Window entries.
func (j *Jepa) Predict() []float64 {
	entries := j.DbIn.Entries
	if len(entries) == 0 {
		return nil
	}
	w := j.Window
	if w <= 0 {
		w = 3
	}
	start := len(entries) - w
	if start < 0 {
		start = 0
	}
	dims := len(entries[0].Data)
	pred := make([]float64, dims)
	count := 0
	for i := start; i < len(entries); i++ {
		for d := 0; d < dims; d++ {
			pred[d] += entries[i].Data[d]
		}
		count++
	}
	for d := 0; d < dims; d++ {
		pred[d] /= float64(count)
	}
	j.last = pred
	return pred
}

// RecordPrediction compares predicted vs actual, records in DbOut, returns surprise.
func (j *Jepa) RecordPrediction(predicted, actual []float64) float64 {
	s := Surprise(predicted, actual)
	cp := make([]float64, len(actual))
	copy(cp, actual)
	j.DbOut.Entries = append(j.DbOut.Entries, Entry{
		Data:     cp,
		Surprise: s,
	})
	return s
}

// Surprise computes L2 distance between two vectors.
func Surprise(a, b []float64) float64 {
	sum := 0.0
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

// CheckConservation verifies that total entries in DbIn ≈ DbOut within tolerance.
func (j *Jepa) CheckConservation(tolerance int) bool {
	diff := len(j.DbIn.Entries) - len(j.DbOut.Entries)
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// GC removes entries from DbOut with surprise below threshold.
func (j *Jepa) GC(threshold float64) {
	filtered := make([]Entry, 0, len(j.DbOut.Entries))
	for _, e := range j.DbOut.Entries {
		if e.Surprise >= threshold {
			filtered = append(filtered, e)
		}
	}
	j.DbOut.Entries = filtered
}
