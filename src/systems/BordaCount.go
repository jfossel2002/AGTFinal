package voting_systems

import (
	"math"
	"sort"
)

func CalculateBordaWinner(candidates []Candidate, voters []Voter) Candidate {
	for i := range candidates {
		candidates[i].NumVotes = 0
	}
	candidates = PerformPairwiseComparisons(candidates, voters)
	for _, voter := range voters {
		// Create a slice to hold distances and corresponding candidate names
		type candidateDistance struct {
			Index    int
			Distance float64
		}

		distances := make([]candidateDistance, len(candidates))

		for i, candidate := range candidates {
			// Calculate the absolute distance from the voter to each candidate
			distances[i] = candidateDistance{
				Index:    i,
				Distance: math.Abs(candidate.Position - voter.Position),
			}
		}

		// Sort candidates by distance from the voter, closest first
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].Distance < distances[j].Distance
		})

		// Distribute points based on distance, closer candidates get more points
		for i, dist := range distances {
			points := (len(candidates) - i - 1) * voter.Number
			candidates[dist.Index].NumVotes += points
		}
	}

	// Find and return the candidate with the most votes
	var winner Candidate
	maxVotes := -1
	for _, candidate := range candidates {
		if candidate.NumVotes > maxVotes {
			maxVotes = candidate.NumVotes
			winner = candidate
		}
	}

	return winner
}
