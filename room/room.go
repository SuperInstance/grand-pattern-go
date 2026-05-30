package room

import (
	"github.com/SuperInstance/grand-pattern-go/jepa"
	"github.com/SuperInstance/grand-pattern-go/murmur"
	"github.com/SuperInstance/grand-pattern-go/vibe"
)

// Room is a single cell in the graph, composing vibe + jepa + murmur.
type Room struct {
	Name    string
	Vibe    *vibe.Vibe
	Jepa    *jepa.Jepa
	Murmurs []*murmur.Murmur
}

// NewRoom creates a room with default state.
func NewRoom(name string) *Room {
	return &Room{
		Name: name,
		Vibe: vibe.NewVibe(),
		Jepa: &jepa.Jepa{Window: 3},
	}
}

// Perceive records a vibe observation into jepa.
func (r *Room) Perceive(ts uint64) {
	r.Jepa.Perceive(r.Vibe.Dims[:], ts)
}

// Tick advances the room one step.
func (r *Room) Tick(ts uint64) {
	r.Perceive(ts)
	// Decay murmurs
	alive := make([]*murmur.Murmur, 0, len(r.Murmurs))
	for _, m := range r.Murmurs {
		murmur.Decay(m)
		if m.TTL > 0 {
			alive = append(alive, m)
		}
	}
	r.Murmurs = alive
}

// AddMurmur adds a murmur to the room.
func (r *Room) AddMurmur(m *murmur.Murmur) {
	r.Murmurs = append(r.Murmurs, m)
}
