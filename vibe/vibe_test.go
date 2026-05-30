package vibe

import (
	"math"
	"testing"
)

func TestNewVibe(t *testing.T) {
	v := NewVibe()
	for i, d := range v.Dims {
		if d != 0 {
			t.Fatalf("dim %d = %f, want 0", i, d)
		}
	}
}

func TestBlend(t *testing.T) {
	a := NewVibe()
	a.Dims[0] = 10
	b := NewVibe()
	b.Dims[0] = 0
	c := a.Blend(b, 0.5)
	if c.Dims[0] != 5 {
		t.Fatalf("blended dim = %f, want 5", c.Dims[0])
	}
}

func TestDistance(t *testing.T) {
	a := NewVibe()
	b := NewVibe()
	a.Dims[0] = 3
	b.Dims[0] = 7
	d := a.Distance(b)
	if math.Abs(d-4) > 1e-9 {
		t.Fatalf("distance = %f, want 4", d)
	}
}

func TestDiffuse(t *testing.T) {
	center := NewVibe()
	center.Dims[0] = 10
	n := NewVibe()
	n.Dims[0] = 0
	result := center.Diffuse([]*Vibe{n}, []float64{1.0}, 0.5)
	if result.Dims[0] != 5 {
		t.Fatalf("diffused = %f, want 5", result.Dims[0])
	}
}

func TestDescription(t *testing.T) {
	v := NewVibe()
	desc := v.QualitativeDescription()
	if desc == "" {
		t.Fatal("empty description")
	}
}
