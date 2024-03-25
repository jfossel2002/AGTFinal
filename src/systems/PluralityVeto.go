package voting_systems

import "math"

func InitiatePluralityVeto(candidates []Candidate, voters []Voter) (Candidate, []Candidate, [][]Candidate) {
	for i := range candidates {
		candidates[i].NumVotes = 0
	}
	CanidateRounds := [][]Candidate{}
	_, winner, candidates, CanidateRounds := SimulatePluralityVeto(candidates, voters, CanidateRounds)
	/*for i := range CanidateRounds {
		fmt.Println("Round ", i)
		fmt.Println(CanidateRounds[i])
	}*/
	return winner, candidates, CanidateRounds
}

func SimulatePluralityVeto(candidates []Candidate, voters []Voter, CanidateRounds [][]Candidate) (float64, Candidate, []Candidate, [][]Candidate) {
	//Give all votes out
	candidates = PluralityVote(voters, candidates)
	CanidateRounds = append(CanidateRounds, append([]Candidate(nil), candidates...))
	//Revoke votes from everyones least Favroite
	//Loop while a canidate has more than 0 votes
	for CheckCanidateVotes(candidates) {
		candidates, CanidateRounds = VetoVote(voters, candidates, CanidateRounds) //Remove the votes from the least favorite
	}

	return float64(candidates[0].NumVotes), candidates[0], candidates, CanidateRounds
}

func checkVotes(voters []Voter) bool {
	for i := 0; i < len(voters); i++ {
		if voters[i].Number > 0 {
			return true
		}
	}
	return false
}

func VetoVote(voters []Voter, candidates []Candidate, CanidateRounds [][]Candidate) ([]Candidate, [][]Candidate) {
	votersCopy := make([]Voter, len(voters))
	copy(votersCopy, voters)
	//For each voter determine the furthest candidate and save it
	//While there are votes to be removed from the voters
	//Pick a random voter, find the furthest candidate and remove 1 vote
	//If the candidate has no votes left, remove the candidate
	for checkVotes(votersCopy) {
		maxDistance := 0.0
		voterPosition := 0
		candidatePosition := 0
		//Pick a random voter whos votes > 0
		for i := 0; i < len(votersCopy); i++ {
			if votersCopy[i].Number > 0 {
				voterPosition = i
				break
			}
		}
		//Find the furthest candidate
		for j := 0; j < len(candidates); j++ {
			distance := math.Abs(votersCopy[voterPosition].Position - candidates[j].Position)
			if distance > maxDistance {
				maxDistance = distance
				candidatePosition = j
			}
		}
		//remove that number of voters from the candidate
		candidates[candidatePosition].NumVotes -= 1
		//fmt.Println(candidates)
		//remove that number of votes from the voter
		votersCopy[voterPosition].Number -= 1
		//Check if the candidate has negative votes
		//If the candidate has no votes left, remove the candidate
		candidates, CanidateRounds = removeNegativeCanidates(candidates, CanidateRounds)

		if len(candidates) == 1 {
			return candidates, CanidateRounds
		}
	}

	return candidates, CanidateRounds
}

// Remove the most negative first
func removeNegativeCanidates(candidates []Candidate, CanidateRounds [][]Candidate) ([]Candidate, [][]Candidate) {
	for i := 0; i < len(candidates); i++ {
		if candidates[i].NumVotes <= 0 {
			//Remove Canidate
			candidates = append(candidates[:i], candidates[i+1:]...)
			CanidateRounds = append(CanidateRounds, append([]Candidate(nil), candidates...)) //Add the round to the list
		}
	}
	return candidates, CanidateRounds
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
