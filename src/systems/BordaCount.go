package voting_systems

/*
 * The Borda Count method is a voting system in which voters rank candidates in order of preference.
 * Points are awarded to each candidate based on their rank, with the top-ranked candidate receiving the most points.
 * The candidate with the most points is declared the winner.
 * This file contains the implementation of the Borda Count method.
 */
import (
	"math"
	"sort"
)

// Function to calculate the winner of an election using the Borda Count method
// Takes in a slice of candidates and a slice of voters
// Returns the winning candidate and a slice of all candidates with their vote counts
func CalculateBordaWinner(candidates []Candidate, voters []Voter) (Candidate, []Candidate) {
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
		if candidate.NumVotes >= maxVotes {
			if candidate.NumVotes == maxVotes {
				//If there is a tie, the canidate that results in the highest distortion wins
				Possible_Winner_Cost := GetSocailCost(candidate, voters)
				Current_Winner_Cost := GetSocailCost(winner, voters)
				opt_cost, _ := DetermineOptimalCanidate(candidates, voters)
				Possbile_Winner_Distortion := GetDistortion(Possible_Winner_Cost, opt_cost)
				Current_Winner_Distortion := GetDistortion(Current_Winner_Cost, opt_cost)
				if Possbile_Winner_Distortion > Current_Winner_Distortion {
					winner = candidate
					maxVotes = candidate.NumVotes

				}
			} else {
				maxVotes = candidate.NumVotes
				winner = candidate
			}
		}
	}

	return winner, candidates
}
