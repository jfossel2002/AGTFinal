package voting_systems

//"fmt"
//"sort"

func CalculateBordaWinner(candidates []Candidate, voters []Voter) Candidate {
	scores := make(map[string]int)
	for _, candidate := range candidates {
		scores[candidate.Name] = 0
	}

	for _, voter := range voters {
		for i, candidate := range candidates {
			scores[candidate.Name] += (len(candidates) - i - 1) * voter.Number
		}
	}
	var winner Candidate
	maxScore := -1
	for _, candidate := range candidates {
		if scores[candidate.Name] > maxScore {
			maxScore = scores[candidate.Name]
			winner = candidate
		}
	}
	return winner
}
