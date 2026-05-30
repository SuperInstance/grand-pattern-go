package murmur

import "math"

// Murmur represents a piece of gossip floating through the network.
type Murmur struct {
	Origin  string
	Level   int
	TTL     int
	Payload float64
}

// MurmurPacket is a serializable gossip message.
type MurmurPacket struct {
	Murmur Murmur
	From   string
}

// GossipRound executes one round of gossip dissemination.
type GossipRound struct {
	Packets []MurmurPacket
}

// Levels returns the number of gossip fanout levels remaining.
func Levels(ttl int) int {
	if ttl <= 0 {
		return 0
	}
	return int(math.Ceil(math.Log2(float64(ttl + 1))))
}

// Decay reduces TTL and level of a murmur.
func Decay(m *Murmur) {
	m.TTL--
	if m.TTL < 0 {
		m.TTL = 0
	}
	m.Level = Levels(m.TTL)
}

// NewMurmur creates a murmur with given origin, payload, and TTL.
func NewMurmur(origin string, payload float64, ttl int) *Murmur {
	return &Murmur{
		Origin:  origin,
		Level:   Levels(ttl),
		TTL:     ttl,
		Payload: payload,
	}
}
