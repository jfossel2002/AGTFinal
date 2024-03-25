package voting_systems

/*
* The Single Transferable Vote (STV) method is a voting system that determines the winner of an election by transferring votes from losing candidates to other candidates based on voter preferences.
* This file contains the implementation of the Single Transferable Vote (STV) method.
 */
import (
	"math"
)

// Function to calculate the winner of an election using the Single Transferable Vote (STV) method
// Takes in a slice of candidates and a slice of voters
// Returns the winning candidate and a slice of all candidates with their vote counts
func InitiateSTV(candidates []Candidate, voters []Voter) (Candidate, []Candidate, [][]Candidate) {
	//Array of rounds of canidate changes
	CanidateRounds := [][]Candidate{}

	for i := range candidates {
		candidates[i].NumVotes = 0
	}
	totalVoters := countTotalVotes(voters)
	candidates = Round(voters, candidates)
	CanidateRounds = append(CanidateRounds, append([]Candidate(nil), candidates...))
	_, winner, candidates, CanidateRounds := SimulateSTV(candidates, totalVoters, voters, CanidateRounds)
	return winner, candidates, CanidateRounds
}

// Count the total number of votes in the election
func countTotalVotes(voters []Voter) int {
	totalVotes := 0
	for _, voter := range voters {
		totalVotes += voter.Number
	}
	return totalVotes
}

// Function to simulate the Single Transferable Vote (STV) voting system
func SimulateSTV(candidates []Candidate, totalVoters int, voters []Voter, CanidateRounds [][]Candidate) (float64, Candidate, []Candidate, [][]Candidate) {
	cost := -1.0
	//Base case
	for i := range candidates {
		if candidates[i].NumVotes > totalVoters/2 {
			return GetSocailCost(candidates[i], voters), candidates[i], candidates, CanidateRounds
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

		candidates = append(candidates[:candidatePosition], candidates[candidatePosition+1:]...)
		candidates = Round(voters, candidates)
		CanidateRounds = append(CanidateRounds, append([]Candidate(nil), candidates...))

		// Recursive call
		cost, winner, candidates, CanidateRounds := SimulateSTV(candidates, totalVoters, voters, CanidateRounds)
		return cost, winner, candidates, CanidateRounds
	}

	//Return empty candidate if no winner is found
	return cost, Candidate{}, candidates, CanidateRounds
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
