package signal

import "testing"

func TestRouterDirect(t *testing.T) {
	r := NewRouter()
	r.AddRoute("a", "out", "b", "in", 0)
	r.Send(Signal{Name: "x", Value: 5, Source: "a"})
	delivered := r.Deliver()
	if len(delivered["b"]) != 1 || delivered["b"][0].Value != 5 {
		t.Fatalf("delivered = %v, want value 5", delivered["b"])
	}
}

func TestDeadband(t *testing.T) {
	r := NewRouter()
	r.AddRoute("a", "out", "b", "in", 0.5)
	r.Send(Signal{Name: "x", Value: 0.3, Source: "a"})
	delivered := r.Deliver()
	if len(delivered["b"]) != 0 {
		t.Fatalf("deadband should filter small signals, got %v", delivered["b"])
	}
}
