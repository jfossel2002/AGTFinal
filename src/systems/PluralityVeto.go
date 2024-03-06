package voting_systems

import "math"

func InitiatePluralityVeto(candidates []Candidate, voters []Voter) Candidate {
	_, winner := SimulatePluralityVeto(candidates, voters)
	return winner
}

func SimulatePluralityVeto(candidates []Candidate, voters []Voter) (float64, Candidate) {
	//Give all votes out
	candidates = PluralityVote(voters, candidates)
	//Revoke votes from everyones least Favroite
	//Loop while a canidate has more than 0 votes
	for CheckCanidateVotes(candidates) {
		candidates = VetoVote(voters, candidates)
	}
	return float64(candidates[0].NumVotes), candidates[0]
}

func VetoVote(voters []Voter, candidates []Candidate) []Candidate {
	//For each voter determine the furthest candidate and save it
	for i := 0; i < len(voters); i++ {
		maxDistance := 0.0
		candidatePosition := 0
		for j := 0; j < len(candidates); j++ {
			distance := math.Abs(voters[i].Position - candidates[j].Position)
			if distance > maxDistance {
				maxDistance = distance
				candidatePosition = j
			}
		}
		//remove that number of voters from the candidate
		candidates[candidatePosition].NumVotes -= voters[i].Number
		//Check if the candidate has negative votes
		if candidates[candidatePosition].NumVotes <= 0 {
			//Remove Canidate
			candidates = append(candidates[:candidatePosition], candidates[candidatePosition+1:]...)
			//Check if there is only one candidate left
			if len(candidates) == 1 {
				return candidates
			}
		}

	}
	return candidates
}

func CheckCanidateVotes(candidates []Candidate) bool {
	if len(candidates) == 1 {
		return false
	}
	for i := 0; i < len(candidates); i++ {
		if candidates[i].NumVotes > 0 {
			return true
		}
	}
	return false
}
