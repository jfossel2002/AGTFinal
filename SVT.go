package main

import "math"

func determineRoundWinner(candidates []Candidate, totalVoters int, voters []Voter) float64 {
	cost := -1.0
	//Base case
	for i := range candidates {
		if candidates[i].numVotes > totalVoters/2 {
			return getSocailCost(candidates[i], voters)
		}
	}

	//Recursive Case
	if cost == -1.0 {
		minVotes := totalVoters // Assuming totalVoters is higher than any candidate's numVotes can be
		candidatePosition := 0
		for j, candidate := range candidates {
			if candidate.numVotes < minVotes {
				minVotes = candidate.numVotes
				candidatePosition = j
			}
		}

		// Remove the candidate with the least votes
		candidates = append(candidates[:candidatePosition], candidates[candidatePosition+1:]...)

		// Redetermine votes based on voters' next preferences
		//printCanidates(candidates)
		candidates = round(voters, candidates)

		// Recursive call to handle the next round
		return determineRoundWinner(candidates, totalVoters, voters)
	}

	return cost
}

// Finds the optimal canidate based on the social cost
func determineOptimal(Candidates []Candidate, voters []Voter) (float64, int) {
	minCost := 1000000.0
	canidatePosition := 0
	for j := 0; j < len(Candidates); j++ {
		getSocailCost(Candidates[j], voters)
		distance := getSocailCost(Candidates[j], voters)
		if distance < minCost {
			minCost = distance
			canidatePosition = j
		}
	}
	//fmt.Println("The optimal canidate is ", Candidates[canidatePosition].name, " ", Candidates[canidatePosition].position, " with a cost of ", minCost)
	return minCost, canidatePosition

}

// Simulates a round by determining the number of votes each canidate gets
func round(voters []Voter, canidates []Candidate) []Candidate {
	//Reset the number of votes for each candidate
	for i := 0; i < len(canidates); i++ {
		canidates[i].numVotes = 0
	}
	//For each voter determine the closest candidate and save it
	for i := 0; i < len(voters); i++ {
		minDistance := 10000.0
		canidatePosition := 0
		for j := 0; j < len(canidates); j++ {
			distance := math.Abs(voters[i].position - canidates[j].position)
			if distance < minDistance {
				minDistance = distance
				canidatePosition = j
			}
		}
		//add that number of voters to the candidate
		canidates[canidatePosition].numVotes += voters[i].number

	}
	return canidates
}
