package vibe

import (
	"math"
	"strings"
)

// Vibe represents a 16-dimensional state vector.
type Vibe struct {
	Dims [16]float64
}

// NewVibe returns a zero-initialized Vibe.
func NewVibe() *Vibe {
	return &Vibe{}
}

// Blend returns a new Vibe that is (1-ratio)*v + ratio*other.
func (v *Vibe) Blend(other *Vibe, ratio float64) *Vibe {
	result := &Vibe{}
	for i := 0; i < 16; i++ {
		result.Dims[i] = v.Dims[i]*(1-ratio) + other.Dims[i]*ratio
	}
	return result
}

// Distance returns Euclidean distance between two Vibes.
func (v *Vibe) Distance(other *Vibe) float64 {
	sum := 0.0
	for i := 0; i < 16; i++ {
		d := v.Dims[i] - other.Dims[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

// Diffuse blends v with weighted neighbors using diffusion coefficient.
func (v *Vibe) Diffuse(neighbors []*Vibe, weights []float64, coeff float64) *Vibe {
	result := &Vibe{}
	copy(result.Dims[:], v.Dims[:])
	if len(neighbors) == 0 || len(weights) == 0 {
		return result
	}
	weightSum := 0.0
	for _, w := range weights {
		weightSum += w
	}
	if weightSum == 0 {
		return result
	}
	for i := 0; i < 16; i++ {
		blended := 0.0
		for j, n := range neighbors {
			if j < len(weights) {
				blended += weights[j] * n.Dims[i]
			}
		}
		blended /= weightSum
		result.Dims[i] = v.Dims[i]*(1-coeff) + blended*coeff
	}
	return result
}

// QualitativeDescription returns a human-readable description of the vibe.
func (v *Vibe) QualitativeDescription() string {
	energy := v.Energy()
	labels := []string{
		"stillness", "calm", "gentle", "steady",
		"lively", "charged", "intense", "frenetic",
	}
	idx := int(math.Min(float64(len(labels)-1), energy/10.0*float64(len(labels))))
	if idx < 0 {
		idx = 0
	}
	// Also describe dominant axis
	maxIdx := 0
	maxVal := 0.0
	for i, d := range v.Dims {
		if math.Abs(d) > maxVal {
			maxVal = math.Abs(d)
			maxIdx = i
		}
	}
	axisNames := []string{
		"harmonic", "rhythmic", "thermal", "luminous",
		"textural", "spatial", "temporal", "gravitic",
		"resonant", "spectral", "kinetic", "magnetic",
		"elastic", "acoustic", "chromatic", "phase",
	}
	return strings.TrimSpace(labels[idx] + " " + axisNames[maxIdx])
}

// FromDescription sets vibe dims heuristically from a text description.
func (v *Vibe) FromDescription(desc string) *Vibe {
	desc = strings.ToLower(desc)
	energyWords := map[string]int{
		"still": 0, "calm": 1, "gentle": 2, "steady": 3,
		"lively": 4, "charged": 5, "intense": 6, "frenetic": 7,
	}
	for word, level := range energyWords {
		if strings.Contains(desc, word) {
			for i := 0; i < 16; i++ {
				v.Dims[i] = float64(level) / 7.0 * 10.0
			}
			break
		}
	}
	axisWords := map[string]int{
		"harmonic": 0, "rhythmic": 1, "thermal": 2, "luminous": 3,
		"textural": 4, "spatial": 5, "temporal": 6, "gravitic": 7,
		"resonant": 8, "spectral": 9, "kinetic": 10, "magnetic": 11,
		"elastic": 12, "acoustic": 13, "chromatic": 14, "phase": 15,
	}
	for word, axis := range axisWords {
		if strings.Contains(desc, word) {
			v.Dims[axis] += 3.0
		}
	}
	return v.Bounded()
}

// GrooveLock returns alignment between two vibes (1.0 = perfectly aligned, 0.0 = orthogonal).
func (v *Vibe) GrooveLock(other *Vibe) float64 {
	dot := 0.0
	normA := 0.0
	normB := 0.0
	for i := 0; i < 16; i++ {
		dot += v.Dims[i] * other.Dims[i]
		normA += v.Dims[i] * v.Dims[i]
		normB += other.Dims[i] * other.Dims[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	cosSim := dot / (math.Sqrt(normA) * math.Sqrt(normB))
	// Normalize from [-1,1] to [0,1]
	return (cosSim + 1) / 2
}

// Energy returns the L2 norm of the vibe vector.
func (v *Vibe) Energy() float64 {
	sum := 0.0
	for _, d := range v.Dims {
		sum += d * d
	}
	return math.Sqrt(sum)
}

// Bounded clamps all dimensions to [0, 10].
func (v *Vibe) Bounded() *Vibe {
	for i := 0; i < 16; i++ {
		v.Dims[i] = math.Max(0, math.Min(10, v.Dims[i]))
	}
	return v
}
