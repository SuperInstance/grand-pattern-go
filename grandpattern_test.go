package grandpattern

import (
	"testing"
)

func TestTickUpdatesPerception(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	if len(r.PerceptionDB) != 1 {
		t.Errorf("expected perception_count=1, got %d", len(r.PerceptionDB))
	}
	if len(r.PredictionDB) != 1 {
		t.Errorf("expected prediction_count=1, got %d", len(r.PredictionDB))
	}
}

func TestPredictGeneratesEmbedding(t *testing.T) {
	r := NewRoom("test")
	pred := r.Predict()
	if pred.Norm() > 1e-9 {
		t.Errorf("first prediction should be near zero, got %f", pred.Norm())
	}
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	r.Tick(2.0, 1, EmbeddingFromValue(2.0))
	pred2 := r.Predict()
	if pred2.Norm() == 0 {
		t.Error("prediction after ticks should be non-zero")
	}
}

func TestBalanceCheckPassesEqual(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	r.Tick(2.0, 1, EmbeddingFromValue(2.0))
	if !r.BalanceCheck() {
		t.Error("balance should pass with equal DB sizes")
	}
}

func TestBalanceCheckFailsUnequal(t *testing.T) {
	r := NewRoom("test")
	r.PerceptionDB = append(r.PerceptionDB, TaggedEmbedding{
		Value: EmbeddingFromValue(1.0), Strength: 1.0, Timestamp: 1.0,
	})
	if r.BalanceCheck() {
		t.Error("balance should fail with unequal DB sizes")
	}
}

func TestVibeComputation(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	if r.Vibe.Strength <= 0 {
		t.Error("vibe strength should be > 0 after tick")
	}
	if r.Vibe.Position.Norm() <= 0 {
		t.Error("vibe position should be non-zero")
	}
}

func TestMergeReducesCount(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	r.Tick(2.0, 1, EmbeddingFromValue(1.01))
	before := len(r.PerceptionDB)
	GC(r, 0.5, 0.99, 0.01)
	if len(r.PerceptionDB) > before {
		t.Error("merge should not increase count")
	}
}

func TestDecayReducesStrengths(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	before := r.PerceptionDB[0].Strength
	GC(r, 999.0, 0.5, 0.0)
	if r.PerceptionDB[0].Strength >= before {
		t.Error("decay should reduce strength")
	}
}

func TestPruneRemovesWeak(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	r.PerceptionDB[0].Strength = 0.001
	if len(r.PredictionDB) > 0 {
		r.PredictionDB[0].Strength = 0.001
	}
	GC(r, 999.0, 1.0, 0.01)
	if len(r.PerceptionDB) != 0 {
		t.Errorf("prune should remove weak embeddings, got %d", len(r.PerceptionDB))
	}
}

func TestFullGCCycle(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1.0, 1, EmbeddingFromValue(1.0))
	r.Tick(2.0, 1, EmbeddingFromValue(1.01))
	r.Tick(3.0, 1, EmbeddingFromValue(5.0))
	report := GC(r, 0.5, 0.9, 0.01)
	if report.Decayed == 0 {
		t.Error("GC should report decayed items")
	}
}

func TestCrossRoomCorrelation(t *testing.T) {
	a := NewRoom("a")
	b := NewRoom("b")
	a.Tick(1.0, 1, EmbeddingFromValue(1.0))
	b.Tick(1.0, 1, EmbeddingFromValue(1.0))
	corr := CrossRoomCorrelation(a, b)
	if corr <= 0.99 {
		t.Errorf("identical embeddings should have correlation ~1.0, got %f", corr)
	}
}

func TestMurmurSendsVibe(t *testing.T) {
	from := NewRoom("from")
	to := NewRoom("to")
	from.Tick(1.0, 1, EmbeddingFromValue(3.0))
	before := len(to.PerceptionDB)
	Murmur(from, to)
	if len(to.PerceptionDB) <= before {
		t.Error("murmur should add to target perception")
	}
}

func TestGraphConstruction(t *testing.T) {
	g := NewGraph()
	g.AddRoom(NewRoom("r1"))
	g.AddRoom(NewRoom("r2"))
	g.AddEdge("r1", "r2", 0)
	if len(g.Rooms) != 2 {
		t.Errorf("expected 2 rooms, got %d", len(g.Rooms))
	}
	if len(g.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(g.Edges))
	}
}

func TestTickThroughGraph(t *testing.T) {
	g := NewGraph()
	g.AddRoom(NewRoom("r1"))
	g.AddRoom(NewRoom("r2"))
	g.AddEdge("r1", "r2", 0)
	g.TickThroughGraph(1.0, 1, EmbeddingFromValue(1.0))
	r1 := g.FindRoom("r1")
	r2 := g.FindRoom("r2")
	if len(r1.PerceptionDB) < 1 {
		t.Error("r1 should have perceptions")
	}
	if len(r2.PerceptionDB) < 2 {
		t.Errorf("r2 should have perceptions + murmur, got %d", len(r2.PerceptionDB))
	}
}
