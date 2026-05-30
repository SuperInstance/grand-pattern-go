package graph

import (
	"math"
	"testing"

	"github.com/SuperInstance/grand-pattern-go/murmur"
	"github.com/SuperInstance/grand-pattern-go/vibe"
)

func TestCellGraphTick(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	g.AddRoom("b")
	g.AddEdge("a", "b")
	g.Tick()
	if g.TickCount != 1 {
		t.Fatalf("tick count = %d, want 1", g.TickCount)
	}
}

func TestCellGraphGossip(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	g.AddRoom("b")
	g.AddEdge("a", "b")
	g.Rooms["a"].AddMurmur(murmur.NewMurmur("a", 1.0, 3))
	g.Gossip()
	if len(g.Rooms["b"].Murmurs) != 1 {
		t.Fatalf("b murmurs = %d, want 1", len(g.Rooms["b"].Murmurs))
	}
	if g.Rooms["b"].Murmurs[0].TTL != 2 {
		t.Fatalf("relayed TTL = %d, want 2", g.Rooms["b"].Murmurs[0].TTL)
	}
}

func TestChainTopology(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	g.AddRoom("b")
	g.AddRoom("c")
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")
	g.SetupDefaultRoutes()

	g.Rooms["a"].Vibe.Dims[0] = 10
	g.Tick()
	g.RouteSignals()

	b := g.Rooms["b"]
	if b.Vibe.Dims[0] < 0.5 {
		t.Fatalf("b dim[0] = %f, expected > 0.5 after routing", b.Vibe.Dims[0])
	}
}

func TestStarTopology(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("center")
	for _, name := range []string{"s1", "s2", "s3"} {
		g.AddRoom(name)
		g.AddEdge("center", name)
	}
	g.SetupDefaultRoutes()
	g.Rooms["center"].Vibe.Dims[0] = 8
	g.Tick()
	g.RouteSignals()

	for _, name := range []string{"s1", "s2", "s3"} {
		if g.Rooms[name].Vibe.Dims[0] < 0.5 {
			t.Fatalf("%s dim[0] = %f, expected > 0.5", name, g.Rooms[name].Vibe.Dims[0])
		}
	}
}

func TestMeshTopology(t *testing.T) {
	g := NewCellGraph(120)
	names := []string{"x", "y", "z"}
	for _, n := range names {
		g.AddRoom(n)
	}
	for _, a := range names {
		for _, b := range names {
			if a != b {
				g.AddEdge(a, b)
			}
		}
	}
	g.SetupDefaultRoutes()
	g.Rooms["x"].Vibe.Dims[1] = 5
	g.Tick()
	g.RouteSignals()

	if g.Rooms["y"].Vibe.Dims[1] < 0.1 {
		t.Fatalf("y dim[1] = %f, expected > 0.1 in mesh", g.Rooms["y"].Vibe.Dims[1])
	}
}

func TestFleetVibe(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	g.AddRoom("b")
	g.Rooms["a"].Vibe.Dims[0] = 4
	g.Rooms["b"].Vibe.Dims[0] = 6
	fv := g.FleetVibe()
	if math.Abs(fv.Dims[0]-5) > 1e-9 {
		t.Fatalf("fleet vibe dim[0] = %f, want 5", fv.Dims[0])
	}
}

func TestDetectAnomaly(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("normal")
	g.AddRoom("hot")
	g.Rooms["hot"].Vibe.Dims[0] = 10
	g.Rooms["hot"].Vibe.Dims[1] = 10
	anomalous := g.DetectAnomaly(5.0)
	if len(anomalous) != 1 || anomalous[0].Name != "hot" {
		t.Fatalf("anomaly = %v, want [hot]", anomalous)
	}
}

func TestSummary(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	s := g.Summary()
	if len(s) == 0 {
		t.Fatal("empty summary")
	}
}

func TestConservationReport(t *testing.T) {
	g := NewCellGraph(120)
	g.AddRoom("a")
	report := g.ConservationReport()
	if _, ok := report["a"]; !ok {
		t.Fatal("missing room a in conservation report")
	}
}

// Suppress unused import
var _ = vibe.NewVibe
var _ = murmur.NewMurmur
