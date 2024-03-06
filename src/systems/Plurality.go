package voting_systems

import "math"

func InitiatePlurality(candidates []Candidate, voters []Voter) Candidate {
	_, winner := SimulatePlurality(candidates, voters)
	return winner
}

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
