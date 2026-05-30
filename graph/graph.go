package graph

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/SuperInstance/grand-pattern-go/jepa"
	"github.com/SuperInstance/grand-pattern-go/murmur"
	"github.com/SuperInstance/grand-pattern-go/room"
	"github.com/SuperInstance/grand-pattern-go/signal"
	"github.com/SuperInstance/grand-pattern-go/tick"
	"github.com/SuperInstance/grand-pattern-go/vibe"
)

// CellGraph is the top-level composition of all primitives.
type CellGraph struct {
	Rooms     map[string]*room.Room
	Edges     [][2]string
	TickCount uint64
	BPM       float64
	Router    *signal.Router
}

// NewCellGraph creates a graph with the given BPM.
func NewCellGraph(bpm float64) *CellGraph {
	return &CellGraph{
		Rooms:  make(map[string]*room.Room),
		BPM:    bpm,
		Router: signal.NewRouter(),
	}
}

// AddRoom adds a named room to the graph.
func (g *CellGraph) AddRoom(name string) string {
	g.Rooms[name] = room.NewRoom(name)
	return name
}

// AddEdge adds a directed edge between two rooms.
func (g *CellGraph) AddEdge(from, to string) {
	g.Edges = append(g.Edges, [2]string{from, to})
}

// Tick advances all rooms by one tick.
func (g *CellGraph) Tick() {
	g.TickCount++
	for _, r := range g.Rooms {
		r.Tick(g.TickCount)
	}
}

// Gossip propagates murmurs along edges.
func (g *CellGraph) Gossip() {
	newMurmurs := map[string][]*murmur.Murmur{}
	for _, edge := range g.Edges {
		from, ok := g.Rooms[edge[0]]
		if !ok {
			continue
		}
		to, ok := g.Rooms[edge[1]]
		if !ok {
			continue
		}
		for _, m := range from.Murmurs {
			relay := &murmur.Murmur{
				Origin:  m.Origin,
				Level:   m.Level,
				TTL:     m.TTL - 1,
				Payload: m.Payload,
			}
			if relay.TTL > 0 {
				newMurmurs[to.Name] = append(newMurmurs[to.Name], relay)
			}
		}
	}
	for name, ms := range newMurmurs {
		if r, ok := g.Rooms[name]; ok {
			r.Murmurs = append(r.Murmurs, ms...)
		}
	}
}

// RouteSignals runs signal routing across all edges.
func (g *CellGraph) RouteSignals() {
	// Send current vibes as signals
	for _, r := range g.Rooms {
		for i, d := range r.Vibe.Dims {
			g.Router.Send(signal.Signal{
				Name:   fmt.Sprintf("dim_%d", i),
				Value:  d,
				Source: r.Name,
			})
		}
	}
	delivered := g.Router.Deliver()
	for destName, sigs := range delivered {
		destRoom, ok := g.Rooms[destName]
		if !ok {
			continue
		}
		// Average delivered signals into vibe
		dimSum := make(map[int]float64)
		dimCount := make(map[int]float64)
		for _, s := range sigs {
			var dim int
			fmt.Sscanf(s.Name, "dim_%d", &dim)
			dimSum[dim] += s.Value
			dimCount[dim]++
		}
		newVibe := vibe.NewVibe()
		copy(newVibe.Dims[:], destRoom.Vibe.Dims[:])
		for dim, sum := range dimSum {
			if dim < 16 && dimCount[dim] > 0 {
				avg := sum / dimCount[dim]
				// Blend 10% of incoming signal
				newVibe.Dims[dim] = newVibe.Dims[dim]*0.9 + avg*0.1
			}
		}
		destRoom.Vibe = newVibe.Bounded()
	}
}

// FleetVibe returns the average vibe across all rooms.
func (g *CellGraph) FleetVibe() *vibe.Vibe {
	avg := vibe.NewVibe()
	n := len(g.Rooms)
	if n == 0 {
		return avg
	}
	for _, r := range g.Rooms {
		for i := 0; i < 16; i++ {
			avg.Dims[i] += r.Vibe.Dims[i]
		}
	}
	for i := 0; i < 16; i++ {
		avg.Dims[i] /= float64(n)
	}
	return avg
}

// FleetSurprise returns average surprise across all rooms' jepa outputs.
func (g *CellGraph) FleetSurprise() float64 {
	total := 0.0
	count := 0
	for _, r := range g.Rooms {
		for _, e := range r.Jepa.DbOut.Entries {
			total += e.Surprise
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// DetectAnomaly returns rooms whose vibe energy exceeds threshold.
func (g *CellGraph) DetectAnomaly(threshold float64) []*room.Room {
	var anomalous []*room.Room
	for _, r := range g.Rooms {
		if r.Vibe.Energy() > threshold {
			anomalous = append(anomalous, r)
		}
	}
	return anomalous
}

// ConservationReport checks jepa conservation for each room.
func (g *CellGraph) ConservationReport() map[string]bool {
	report := make(map[string]bool)
	for name, r := range g.Rooms {
		report[name] = r.Jepa.CheckConservation(2)
	}
	return report
}

// Summary returns a human-readable graph summary.
func (g *CellGraph) Summary() string {
	var b strings.Builder
	roomNames := make([]string, 0, len(g.Rooms))
	for name := range g.Rooms {
		roomNames = append(roomNames, name)
	}
	sort.Strings(roomNames)

	fmt.Fprintf(&b, "CellGraph | tick=%d bpm=%.1f rooms=%d edges=%d\n",
		g.TickCount, g.BPM, len(g.Rooms), len(g.Edges))
	for _, name := range roomNames {
		r := g.Rooms[name]
		fv := g.FleetVibe()
		_ = fv
		fmt.Fprintf(&b, "  [%s] energy=%.2f murmur=%d desc=%s\n",
			name, r.Vibe.Energy(), len(r.Murmurs), r.Vibe.QualitativeDescription())
	}
	return b.String()
}

// SetupDefaultRoutes creates routes for all edges with direct algorithm (deadband=0).
func (g *CellGraph) SetupDefaultRoutes() {
	for _, edge := range g.Edges {
		g.Router.AddRoute(edge[0], "out", edge[1], "in", 0)
	}
}

// Compute the L2 norm — helper (avoids importing vibe just for math).
func l2norm(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x * x
	}
	return math.Sqrt(s)
}

// Suppress unused import
var _ = jepa.Surprise
var _ = murmur.Levels
var _ = tick.Tempo
