package routing

import (
	"sort"

	"github.com/nyxora-network/nyxora/internal/transport"
)

type PathScore struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Score     float64 `json:"score"`
	Latency   float64 `json:"latency"`
	Jitter    float64 `json:"jitter"`
	Loss      float64 `json:"packet_loss"`
	Stability float64 `json:"stability"`
	Bandwidth int     `json:"bandwidth"`
}

type Scorer struct {
	weights transport.ScoringWeights
}

var DefaultWeights = transport.ScoringWeights{
	Latency:   0.30,
	Loss:      0.25,
	Jitter:    0.15,
	Stability: 0.20,
}

func NewScorer() *Scorer {
	return &Scorer{weights: DefaultWeights}
}

func NewScorerWithWeights(w transport.ScoringWeights) *Scorer {
	return &Scorer{weights: w}
}

func (s *Scorer) Score(latency, jitter, loss, stability float64, bandwidth int) float64 {
	m := &transport.Metrics{
		LatencyMs:  latency,
		JitterMs:   jitter,
		PacketLoss: loss,
		Stability:  stability,
		Bandwidth:  bandwidth,
	}
	return transport.ComputeScore(m, s.weights)
}

func (s *Scorer) Rank(scores []PathScore) []PathScore {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	return scores
}

func (s *Scorer) Best(scores []PathScore) *PathScore {
	if len(scores) == 0 {
		return nil
	}
	ranked := s.Rank(scores)
	return &ranked[0]
}
