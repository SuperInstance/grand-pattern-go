package room

import (
	"testing"

	"github.com/SuperInstance/grand-pattern-go/murmur"
)

func TestNewRoom(t *testing.T) {
	r := NewRoom("test")
	if r.Name != "test" {
		t.Errorf("expected name 'test', got %q", r.Name)
	}
	if r.Vibe == nil {
		t.Error("expected non-nil Vibe")
	}
	if r.Jepa == nil {
		t.Error("expected non-nil Jepa")
	}
}

func TestRoomTick(t *testing.T) {
	r := NewRoom("test")
	r.Tick(1)
	r.Tick(2)
	if len(r.Jepa.DbIn.Entries) != 2 {
		t.Errorf("expected 2 jepa entries, got %d", len(r.Jepa.DbIn.Entries))
	}
}

func TestRoomMurmurLifecycle(t *testing.T) {
	r := NewRoom("test")
	m := &murmur.Murmur{Origin: "other", TTL: 2}
	r.AddMurmur(m)
	if len(r.Murmurs) != 1 {
		t.Errorf("expected 1 murmur, got %d", len(r.Murmurs))
	}
	r.Tick(1) // decay → TTL 1
	if len(r.Murmurs) != 1 {
		t.Errorf("expected 1 murmur after first tick, got %d", len(r.Murmurs))
	}
	r.Tick(2) // decay → TTL 0 → removed
	if len(r.Murmurs) != 0 {
		t.Errorf("expected 0 murmurs after expiry, got %d", len(r.Murmurs))
	}
}
