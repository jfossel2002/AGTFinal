package voting_systems

/*
* The Plurality method is a voting system that determines the winner of an election by selecting the candidate with the most votes.
* Each voter selects their preferred candidate, and the candidate with the most votes is declared the winner.
* This file contains the implementation of the Plurality method.
 */
import "math"

// Function to calculate the winner of an election using the Plurality method
// Takes in a slice of candidates and a slice of voters
// Returns the winning candidate and a slice of all candidates with their vote counts
func InitiatePlurality(candidates []Candidate, voters []Voter) (Candidate, []Candidate) {
	for i := range candidates {
		candidates[i].NumVotes = 0
	}
	_, winner := SimulatePlurality(candidates, voters)
	return winner, candidates
}

// Function to simulate the Plurality voting system
func SimulatePlurality(candidates []Candidate, voters []Voter) (float64, Candidate) {
	candidates = PluralityVote(voters, candidates)
	//Determine canidate with the most votes
	maxVotes := 0
	canidatePosition := 0
	for i := 0; i < len(candidates); i++ {
		if candidates[i].NumVotes > maxVotes {
			maxVotes = candidates[i].NumVotes
			canidatePosition = i
		}
	}
	//Return the winner
	return float64(candidates[canidatePosition].NumVotes), candidates[canidatePosition]
}

// Simulates a round by determining the number of votes each canidate gets
func PluralityVote(voters []Voter, canidates []Candidate) []Candidate {
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
