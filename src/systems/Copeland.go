package voting_systems

/*
 * The Copeland method is a voting system that determines the winner of an election by comparing each pair of candidates.
 * A candidate receives a point for each other candidate they beat in a pairwise comparison.
 * The candidate with the most points is declared the winner.
 * This file contains the implementation of the Copeland method.
 */
import (
	"math"
)

// Function to determine if a voter prefers one candidate over another
// Takes in a voter and two candidates
// Returns true if the voter prefers candidateA over candidateB
func VoterPrefers(voter Voter, candidateA, candidateB Candidate) bool {
	return math.Abs(voter.Position-candidateA.Position) < math.Abs(voter.Position-candidateB.Position)
}

// Function to perform pairwise comparisons between candidates
// Takes in a slice of candidates and a slice of voters
// Returns a slice of candidates with updated vote counts
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

// Function to determine the winner of an election using the Copeland method
// Takes in a slice of candidates and a slice of voters
// Returns the winning candidate and a slice of all candidates with their vote counts
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
		if candidate.NumVotes >= maxScore {
			if candidate.NumVotes == maxScore {
				//If there is a tie, the canidate that results in the highest distortion wins
				Possible_Winner_Cost := GetSocailCost(candidate, voters)
				Current_Winner_Cost := GetSocailCost(winner, voters)
				opt_cost, _ := DetermineOptimalCanidate(candidates, voters)
				Possbile_Winner_Distortion := GetDistortion(Possible_Winner_Cost, opt_cost)
				Current_Winner_Distortion := GetDistortion(Current_Winner_Cost, opt_cost)
				if Possbile_Winner_Distortion > Current_Winner_Distortion {
					winner = candidate
					maxScore = candidate.NumVotes

				}
			} else {
				maxScore = candidate.NumVotes
				winner = candidate
			}
		}
	}

	return winner, candidates
}
