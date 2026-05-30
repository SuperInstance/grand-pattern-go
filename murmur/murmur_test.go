package murmur

import "testing"

func TestNewMurmur(t *testing.T) {
	m := NewMurmur("room1", 42.0, 5)
	if m.Origin != "room1" || m.Payload != 42.0 || m.TTL != 5 {
		t.Fatalf("murmur = %+v, unexpected", m)
	}
}

func TestDecay(t *testing.T) {
	m := NewMurmur("a", 1.0, 3)
	Decay(m)
	if m.TTL != 2 {
		t.Fatalf("TTL after decay = %d, want 2", m.TTL)
	}
}

func TestGossip(t *testing.T) {
	m := NewMurmur("origin", 3.14, 4)
	pkts := []MurmurPacket{
		{Murmur: *m, From: "a"},
		{Murmur: *m, From: "b"},
	}
	gr := &GossipRound{Packets: pkts}
	if len(gr.Packets) != 2 {
		t.Fatalf("packets = %d, want 2", len(gr.Packets))
	}
}
