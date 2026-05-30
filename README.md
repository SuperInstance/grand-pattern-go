# Grand Pattern - Go Implementation

Fibonacci Dual-Direction Architecture — a cellular graph system with dual perception/prediction databases.

## Architecture

- **Perception DB (Z_in)**: stores incoming sensor embeddings
- **Prediction DB (Z_out)**: stores predicted future embeddings
- **JEPA mapping**: cross-DB comparison, computes prediction error (surprise)
- **Double-entry bookkeeping**: every tick updates BOTH databases, must balance
- **Vibe**: (position, velocity, acceleration) tuple on the embedding manifold
- **GC**: 3-phase (merge similar → decay old → prune weak)
- **Cellular graph**: rooms as nodes, algorithms as edges, murmur as gossip protocol

## Test

```bash
go test -v ./...
```

## Dependencies

None. Pure Go, stdlib only.
