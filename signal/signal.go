package signal

import "math"

// Signal represents a named floating-point value with a source port.
type Signal struct {
	Name   string
	Value  float64
	Source string
}

// Port is a named input/output on a room.
type Port struct {
	Name     string
	Incoming []Signal
}

// Route connects a source room+port to a destination room+port.
type Route struct {
	From     string
	FromPort string
	To       string
	ToPort   string
	Deadband float64
}

// Router manages signal routing between rooms.
type Router struct {
	Routes []Route
	Buffer map[string][]Signal
}

// NewRouter creates an empty router.
func NewRouter() *Router {
	return &Router{Buffer: make(map[string][]Signal)}
}

// AddRoute adds a new route with optional deadband (0 = no deadband).
func (r *Router) AddRoute(from, fromPort, to, toPort string, deadband float64) {
	r.Routes = append(r.Routes, Route{
		From: from, FromPort: fromPort, To: to, ToPort: toPort, Deadband: deadband,
	})
}

// Send enqueues a signal for routing.
func (r *Router) Send(sig Signal) {
	r.Buffer[sig.Source] = append(r.Buffer[sig.Source], sig)
}

// Deliver routes all buffered signals according to routes, applying algorithms.
// Returns map of destination room -> delivered signals.
func (r *Router) Deliver() map[string][]Signal {
	result := make(map[string][]Signal)
	for _, route := range r.Routes {
		srcSignals, ok := r.Buffer[route.From]
		if !ok {
			continue
		}
		for _, sig := range srcSignals {
			// Apply deadband
			if route.Deadband > 0 && math.Abs(sig.Value) < route.Deadband {
				continue
			}
			delivered := Signal{
				Name:   sig.Name,
				Value:  applyAlgorithm(sig.Value, route),
				Source: route.From,
			}
			result[route.To] = append(result[route.To], delivered)
		}
	}
	r.Buffer = make(map[string][]Signal)
	return result
}

// applyAlgorithm applies one of 6 routing algorithms based on deadband encoding.
// deadband ranges: [0, 0.1)=direct, [0.1,0.3)=smooth, [0.3,0.5)=boost,
// [0.5,0.7)=compress, [0.7,0.9)=invert, [0.9+]=quantize
func applyAlgorithm(val float64, route Route) float64 {
	db := route.Deadband
	switch {
	case db < 0.1:
		return val // direct
	case db < 0.3:
		return val * 0.9 // smooth (low-pass approximation)
	case db < 0.5:
		return val * 1.5 // boost
	case db < 0.7:
		if val > 0 {
			return math.Log1p(val) // compress
		}
		return -math.Log1p(-val)
	case db < 0.9:
		return -val // invert
	default:
		return math.Round(val) // quantize
	}
}
