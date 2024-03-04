// File for join functions for scenarios and simulator, i.e. voting systems?
package primary

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

// Prints the voters
func printVoters(voters []Voter) {
	//Sort voters based on position
	sort.Slice(voters, func(i, j int) bool {
		return voters[i].Position < voters[j].Position
	})
	fmt.Println("\nVoters: ")
	for i := 0; i < len(voters); i++ {
		printVoter(voters[i])
	}
}

// Prints a single voter
func printVoter(voter Voter) {
	fmt.Println(voter.Name, " ", voter.Number)
	fmt.Println("  Position ", voter.Position)
}

// Prints the canidates
func printCanidates(canidates []Candidate) {
	fmt.Println("\nCanidates: ")
	//Sort canidates based on position
	sort.Slice(canidates, func(i, j int) bool {
		return canidates[i].Position < canidates[j].Position
	})
	for i := 0; i < len(canidates); i++ {
		printCanidate(canidates[i])
	}
}

// Prints a single canidate
func printCanidate(canidate Candidate) {
	fmt.Println(canidate.Name, " ", canidate.Position)
	fmt.Println("  Votes ", canidate.NumVotes)
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
