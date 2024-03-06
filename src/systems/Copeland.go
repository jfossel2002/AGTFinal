package voting_systems

import (
	"math"
)

func VoterPrefers(voter Voter, candidateA, candidateB Candidate) bool {
	return math.Abs(voter.Position-candidateA.Position) < math.Abs(voter.Position-candidateB.Position)
}

func PerformPairwiseComparisons(candidates []Candidate, voters []Voter) map[string]int {
	scores := make(map[string]int)
	// Initialize scores to zero
	for _, candidate := range candidates {
		scores[candidate.Name] = 0
	}

	// Perform pairwise comparisons
	for i, candidateA := range candidates {
		for j := i + 1; j < len(candidates); j++ {
			candidateB := candidates[j]
			winsA, winsB := 0, 0

			// Count the number of times each candidate wins in pairwise comparisons
			for _, voter := range voters {
				if VoterPrefers(voter, candidateA, candidateB) {
					winsA++
				} else if VoterPrefers(voter, candidateB, candidateA) {
					winsB++
				}
			}

			// Update scores based on pairwise comparisons
			if winsA > winsB {
				scores[candidateA.Name]++
			} else if winsB > winsA {
				scores[candidateB.Name]++
			} else { // Handle ties
				scores[candidateA.Name]++
				scores[candidateB.Name]++
			}
		}
	}

	return scores
}

func CalculateCopelandScore(candidates []Candidate, voters []Voter) map[string]int {
	// Perform pairwise comparisons
	scores := PerformPairwiseComparisons(candidates, voters)

	// Convert scores to Copeland scores (subtracting losses)
	copelandScores := make(map[string]int)
	for candidateName, score := range scores {
		copelandScores[candidateName] = score - (len(candidates) - 1 - score)
	}

	return copelandScores
}

func DetermineCopelandWinner(candidates []Candidate, voters []Voter) Candidate {
	copelandScores := CalculateCopelandScore(candidates, voters)

	// Find candidate with highest Copeland score
	var winner Candidate
	maxScore := -1
	for _, candidate := range candidates {
		if score, ok := copelandScores[candidate.Name]; ok && score > maxScore {
			maxScore = score
			winner = candidate
		}
	}

	return winner
}
