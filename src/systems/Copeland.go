package voting_systems

import (
	"math"
)

func VoterPrefers(voter Voter, candidateA, candidateB Candidate) bool {
	return math.Abs(voter.Position-candidateA.Position) < math.Abs(voter.Position-candidateB.Position)
}

func PerformPairwiseComparisons(candidates []Candidate, voters []Voter) []Candidate {
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
					winsA += voter.Number
				} else if VoterPrefers(voter, candidateB, candidateA) {
					winsB += voter.Number
				}
			}

			// Update scores based on pairwise comparisons
			if winsA > winsB {
				scores[candidateA.Name]++
				candidates[i].NumVotes++
			} else if winsB > winsA {
				scores[candidateB.Name]++
				candidates[j].NumVotes++
			} else { // Handle ties
				scores[candidateA.Name]++
				scores[candidateB.Name]++
				candidates[i].NumVotes++
				candidates[j].NumVotes++
			}
		}
	}

	return candidates
}

func DetermineCopelandWinner(candidates []Candidate, voters []Voter) (Candidate, []Candidate) {
	//Reset all candidate scores
	for i := range candidates {
		candidates[i].NumVotes = 0
	}
	candidates = PerformPairwiseComparisons(candidates, voters)

	// Find candidate with highest Copeland score
	var winner Candidate
	maxScore := -1
	for _, candidate := range candidates {
		if candidate.NumVotes > maxScore {
			maxScore = candidate.NumVotes
			winner = candidate
		}
	}

	return winner, candidates
}
