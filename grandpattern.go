package grandpattern

import "math"

const EmbeddingDim = 8

type Embedding [EmbeddingDim]float64

func EmbeddingZero() Embedding {
	return Embedding{}
}

func EmbeddingFromValue(base float64) Embedding {
	var e Embedding
	for i := 0; i < EmbeddingDim; i++ {
		e[i] = base + float64(i)*0.1
	}
	return e
}

func (a Embedding) Add(b Embedding) Embedding {
	var r Embedding
	for i := 0; i < EmbeddingDim; i++ {
		r[i] = a[i] + b[i]
	}
	return r
}

func (a Embedding) Sub(b Embedding) Embedding {
	var r Embedding
	for i := 0; i < EmbeddingDim; i++ {
		r[i] = a[i] - b[i]
	}
	return r
}

func (a Embedding) Scale(s float64) Embedding {
	var r Embedding
	for i := 0; i < EmbeddingDim; i++ {
		r[i] = a[i] * s
	}
	return r
}

func (a Embedding) Dot(b Embedding) float64 {
	var s float64
	for i := 0; i < EmbeddingDim; i++ {
		s += a[i] * b[i]
	}
	return s
}

func (a Embedding) Norm() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Embedding) CosineSimilarity(b Embedding) float64 {
	d := a.Norm() * b.Norm()
	if d < 1e-12 {
		return 0.0
	}
	return a.Dot(b) / d
}

func (a Embedding) EuclideanDist(b Embedding) float64 {
	var s float64
	for i := 0; i < EmbeddingDim; i++ {
		d := a[i] - b[i]
		s += d * d
	}
	return math.Sqrt(s)
}

func Centroid(db []TaggedEmbedding) Embedding {
	if len(db) == 0 {
		return EmbeddingZero()
	}
	c := EmbeddingZero()
	for _, te := range db {
		c = c.Add(te.Value)
	}
	return c.Scale(1.0 / float64(len(db)))
}

type TaggedEmbedding struct {
	Value     Embedding
	Strength  float64
	Timestamp float64
}

type Vibe struct {
	Position     Embedding
	Velocity     Embedding
	Acceleration Embedding
	Strength     float64
}

type Room struct {
	ID             string
	PerceptionDB   []TaggedEmbedding
	PredictionDB   []TaggedEmbedding
	Vibe           Vibe
	prevPosition   Embedding
	prevVelocity   Embedding
}

func NewRoom(id string) *Room {
	return &Room{
		ID:           id,
		PerceptionDB: nil,
		PredictionDB: nil,
	}
}

func (r *Room) Tick(timestamp float64, sensorID int, perception Embedding) float64 {
	r.PerceptionDB = append(r.PerceptionDB, TaggedEmbedding{
		Value:     perception,
		Strength:  1.0,
		Timestamp: timestamp,
	})

	pred := r.Predict()
	r.PredictionDB = append(r.PredictionDB, TaggedEmbedding{
		Value:     pred,
		Strength:  1.0,
		Timestamp: timestamp,
	})

	err := perception.EuclideanDist(pred)
	r.ComputeVibe()
	return err
}

func (r *Room) Predict() Embedding {
	return r.Vibe.Position.Add(r.Vibe.Velocity)
}

func (r *Room) BalanceCheck() bool {
	return len(r.PerceptionDB) == len(r.PredictionDB)
}

func (r *Room) ComputeVibe() {
	oldVel := r.Vibe.Velocity

	r.Vibe.Position = Centroid(r.PerceptionDB)
	r.Vibe.Velocity = r.Vibe.Position.Sub(r.prevPosition)
	r.Vibe.Acceleration = r.Vibe.Velocity.Sub(r.prevVelocity)
	r.Vibe.Strength = float64(len(r.PerceptionDB))

	r.prevPosition = r.Vibe.Position
	r.prevVelocity = oldVel
}

type GCReport struct {
	Merged  int
	Decayed int
	Pruned  int
}

func mergeSimilar(db *[]TaggedEmbedding, threshold float64) int {
	arr := *db
	n := len(arr)
	removed := make([]bool, n)
	merged := 0

	for i := 0; i < n; i++ {
		if removed[i] {
			continue
		}
		for j := i + 1; j < n; j++ {
			if removed[j] {
				continue
			}
			if arr[i].Value.EuclideanDist(arr[j].Value) < threshold {
				arr[i].Value = arr[i].Value.Add(arr[j].Value).Scale(0.5)
				arr[i].Strength += arr[j].Strength
				removed[j] = true
				merged++
			}
		}
	}

	write := 0
	for i := 0; i < n; i++ {
		if !removed[i] {
			arr[write] = arr[i]
			write++
		}
	}
	*db = arr[:write]
	return merged
}

func applyDecay(db *[]TaggedEmbedding, rate float64) int {
	for i := range *db {
		(*db)[i].Strength *= rate
	}
	return len(*db)
}

func pruneWeak(db *[]TaggedEmbedding, minStrength float64) int {
	before := len(*db)
	write := 0
	for _, te := range *db {
		if te.Strength >= minStrength {
			(*db)[write] = te
			write++
		}
	}
	*db = (*db)[:write]
	return before - write
}

func GC(room *Room, mergeThreshold, decayRate, minStrength float64) GCReport {
	merged := mergeSimilar(&room.PerceptionDB, mergeThreshold) +
		mergeSimilar(&room.PredictionDB, mergeThreshold)
	decayed := applyDecay(&room.PerceptionDB, decayRate) +
		applyDecay(&room.PredictionDB, decayRate)
	pruned := pruneWeak(&room.PerceptionDB, minStrength) +
		pruneWeak(&room.PredictionDB, minStrength)
	return GCReport{Merged: merged, Decayed: decayed, Pruned: pruned}
}

type Edge struct {
	FromID    string
	ToID      string
	Algorithm int
}

type CellularGraph struct {
	Rooms []*Room
	Edges []Edge
}

func NewGraph() *CellularGraph {
	return &CellularGraph{}
}

func (g *CellularGraph) AddRoom(room *Room) {
	g.Rooms = append(g.Rooms, room)
}

func (g *CellularGraph) AddEdge(from, to string, algorithm int) {
	g.Edges = append(g.Edges, Edge{FromID: from, ToID: to, Algorithm: algorithm})
}

func (g *CellularGraph) FindRoom(id string) *Room {
	for _, r := range g.Rooms {
		if r.ID == id {
			return r
		}
	}
	return nil
}

func Murmur(from, to *Room) Embedding {
	vibePos := from.Vibe.Position
	to.PerceptionDB = append(to.PerceptionDB, TaggedEmbedding{
		Value:    vibePos,
		Strength: 0.5,
	})
	return vibePos
}

func CrossRoomCorrelation(a, b *Room) float64 {
	return a.Vibe.Position.CosineSimilarity(b.Vibe.Position)
}

func (g *CellularGraph) TickThroughGraph(timestamp float64, sensorID int, perception Embedding) {
	for _, room := range g.Rooms {
		room.Tick(timestamp, sensorID, perception)
	}
	for _, edge := range g.Edges {
		from := g.FindRoom(edge.FromID)
		to := g.FindRoom(edge.ToID)
		if from != nil && to != nil {
			Murmur(from, to)
		}
	}
}
