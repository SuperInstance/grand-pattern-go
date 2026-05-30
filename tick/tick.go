package tick

// No imports needed

// Tick represents a single discrete time step.
type Tick struct {
	Count    uint64
	BPM      float64
	Elapsed  float64 // seconds
}

// TickSchedule manages scheduled events.
type TickSchedule struct {
	Events map[uint64][]func()
}

// NewTickSchedule creates an empty schedule.
func NewTickSchedule() *TickSchedule {
	return &TickSchedule{Events: make(map[uint64][]func())}
}

// At registers a callback at a specific tick count.
func (ts *TickSchedule) At(tick uint64, fn func()) {
	ts.Events[tick] = append(ts.Events[tick], fn)
}

// Fire executes all callbacks registered for the given tick.
func (ts *TickSchedule) Fire(tick uint64) {
	for _, fn := range ts.Events[tick] {
		fn()
	}
}

// Tempo computes seconds per tick from BPM.
func Tempo(bpm float64) float64 {
	if bpm <= 0 {
		return 1.0
	}
	return 60.0 / bpm
}

// TMinusEvent returns how many ticks remain until a target tick.
func TMinusEvent(current, target uint64) uint64 {
	if target <= current {
		return 0
	}
	return target - current
}

// AdaptBPM smoothly adjusts BPM toward target.
func AdaptBPM(current, target, rate float64) float64 {
	return current + (target-current)*rate
}

// Advance moves a tick forward, returning the new Tick state.
func (t *Tick) Advance() *Tick {
	t.Count++
	t.Elapsed += Tempo(t.BPM)
	return t
}
