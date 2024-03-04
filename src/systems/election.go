// File for join functions for scenarios and simulator, i.e. voting systems?
package voting_systems

import (
	"fmt"
	"math"
	"sort"
)

type Candidate struct {
	Name     string
	Position float64
	NumVotes int
}

type Voter struct {
	Name     string
	Number   int
	Position float64
}

// Returns the distortion of the election
func GetDistortion(winnterCost, optimalCost float64) float64 {
	return (winnterCost / optimalCost)
}

// Returns the social cost of a canidate
func GetSocailCost(c Candidate, voters []Voter) float64 {
	cost := 0.0
	for i := 0; i < len(voters); i++ {
		cost += math.Abs(voters[i].Position-c.Position) * float64(voters[i].Number)
	}
	return cost
}

// Finds the optimal canidate based on the social cost
func DetermineOptimalCanidate(Candidates []Candidate, voters []Voter) (float64, Candidate) {
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
	return minCost, Candidates[canidatePosition]
}

func PrintAllCosts(Candidates []Candidate, voters []Voter) {
	for i := 0; i < len(Candidates); i++ {
		fmt.Println("The cost for ", Candidates[i].Name, " ", Candidates[i].Position, " is ", GetSocailCost(Candidates[i], voters))
	}
}

// Prints the voters
func PrintVoters(voters []Voter) {
	//Sort voters based on position
	sort.Slice(voters, func(i, j int) bool {
		return voters[i].Position < voters[j].Position
	})
	fmt.Println("\nVoters: ")
	for i := 0; i < len(voters); i++ {
		PrintVoter(voters[i])
	}
}

// Prints a single voter
func PrintVoter(voter Voter) {
	fmt.Println(voter.Name, " ", voter.Number)
	fmt.Println("  Position ", voter.Position)
}

// Prints the canidates
func PrintCanidates(canidates []Candidate) {
	fmt.Println("\nCanidates: ")
	//Sort canidates based on position
	sort.Slice(canidates, func(i, j int) bool {
		return canidates[i].Position < canidates[j].Position
	})
	for i := 0; i < len(canidates); i++ {
		PrintCanidate(canidates[i])
	}
}

// Prints a single canidate
func PrintCanidate(canidate Candidate) {
	fmt.Println(canidate.Name, " ", canidate.Position)
	fmt.Println("  Votes ", canidate.NumVotes)
}
