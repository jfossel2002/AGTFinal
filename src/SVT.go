package primary

import (
	"math"
)

func DetermineRoundWinner(candidates []Candidate, totalVoters int, voters []Voter) float64 {
	cost := -1.0
	//Base case
	for i := range candidates {
		if candidates[i].NumVotes > totalVoters/2 {
			return GetSocailCost(candidates[i], voters)
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
		//printCanidates(candidates)
		candidates = Round(voters, candidates)

		// Recursive call to handle the next round
		return DetermineRoundWinner(candidates, totalVoters, voters)
	}

	return cost
}

// Finds the optimal canidate based on the social cost
func DetermineOptimal(Candidates []Candidate, voters []Voter) (float64, int) {
	minCost := 1000000.0
	canidatePosition := 0
	for j := 0; j < len(Candidates); j++ {
		GetSocailCost(Candidates[j], voters)
		distance := GetSocailCost(Candidates[j], voters)
		if distance < minCost {
			minCost = distance
			canidatePosition = j
		}
	}
	//fmt.Println("The optimal canidate is ", Candidates[canidatePosition].name, " ", Candidates[canidatePosition].position, " with a cost of ", minCost)
	return minCost, canidatePosition

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
