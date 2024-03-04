package voting_systems

import (
	"math"
)

func InitiateSTV(candidates []Candidate, totalVoters int, voters []Voter) Candidate {
	Round([]Voter{}, []Candidate{})
	_, winner := SimulateSTV([]Candidate{}, 0, []Voter{})
	return winner
}

func SimulateSTV(candidates []Candidate, totalVoters int, voters []Voter) (float64, Candidate) {
	cost := -1.0
	//Base case
	for i := range candidates {
		if candidates[i].NumVotes > totalVoters/2 {
			return GetSocailCost(candidates[i], voters), candidates[i]
		}
	}

	//Recursive Case
	if cost == -1.0 {
		minVotes := totalVoters // Assuming totalVoters is higher than any candidate's numVotes can be
		candidatePosition := 0
		for j, candidate := range candidates {
			if candidate.NumVotes < minVotes {
				minVotes = candidate.NumVotes
				candidatePosition = j
			}
		}

		// Remove the candidate with the least votes
		candidates = append(candidates[:candidatePosition], candidates[candidatePosition+1:]...)

		// Redetermine votes based on voters' next preferences
		candidates = Round(voters, candidates)

		// Recursive call to handle the next round
		return SimulateSTV(candidates, totalVoters, voters)
	}

	//Return empty candidate if no winner is found
	return cost, Candidate{}
}

// Simulates a round by determining the number of votes each canidate gets
func Round(voters []Voter, canidates []Candidate) []Candidate {
	//Reset the number of votes for each candidate
	for i := 0; i < len(canidates); i++ {
		canidates[i].NumVotes = 0
	}
	//For each voter determine the closest candidate and save it
	for i := 0; i < len(voters); i++ {
		minDistance := 10000.0
		canidatePosition := 0
		for j := 0; j < len(canidates); j++ {
			distance := math.Abs(voters[i].Position - canidates[j].Position)
			if distance < minDistance {
				minDistance = distance
				canidatePosition = j
			}
		}
		//add that number of voters to the candidate
		canidates[canidatePosition].NumVotes += voters[i].Number

	}
	return canidates
}
