package voting_systems

/*
* The Plurality Veto method is a voting system that determines the winner of an election by selecting the candidate with the most votes
* and then removing votes from the least preferred candidates until one candidate has more than half of the votes.
* This file contains the implementation of the Plurality Veto method.
 */
import "math"

// Function to calculate the winner of an election using the Plurality Veto method
// Takes in a slice of candidates and a slice of voters
// Returns the winning candidate and a slice of all candidates with their vote counts
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

// Function to simulate the Plurality Veto voting system
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

// Checks to see if there are any voters with remaining votes
func checkVotes(voters []Voter) bool {
	for i := 0; i < len(voters); i++ {
		if voters[i].Number > 0 {
			return true
		}
	}
	return false
}

// Simulates the veto system
// For each voter determine the furthest candidate and save it
// While there are votes to be removed from the voters
// Pick a random voter, find the furthest candidate and remove 1 vote
// If the candidate has no votes left, remove the candidate
func VetoVote(voters []Voter, candidates []Candidate, CanidateRounds [][]Candidate) ([]Candidate, [][]Candidate) {
	votersCopy := make([]Voter, len(voters))
	copy(votersCopy, voters)
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

// Removes canidates with <= 0 votes
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

// Checks to see if any candidates have > 0 votes or if there is only one left
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
