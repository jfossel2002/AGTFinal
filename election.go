// File for join functions for scenarios and simulator, i.e. voting systems?
package main

import (
	"fmt"
	"math"
	"sort"
)

type Candidate struct {
	name     string
	position float64
	numVotes int
}

type Voter struct {
	name     string
	number   int
	position float64
}

// Prints the voters
func printVoters(voters []Voter) {
	//Sort voters based on position
	sort.Slice(voters, func(i, j int) bool {
		return voters[i].position < voters[j].position
	})
	fmt.Println("\nVoters: ")
	for i := 0; i < len(voters); i++ {
		printVoter(voters[i])
	}
}

// Prints a single voter
func printVoter(voter Voter) {
	fmt.Println(voter.name, " ", voter.number)
	fmt.Println("  Position ", voter.position)
}

// Prints the canidates
func printCanidates(canidates []Candidate) {
	fmt.Println("\nCanidates: ")
	//Sort canidates based on position
	sort.Slice(canidates, func(i, j int) bool {
		return canidates[i].position < canidates[j].position
	})
	for i := 0; i < len(canidates); i++ {
		printCanidate(canidates[i])
	}
}

// Prints a single canidate
func printCanidate(canidate Candidate) {
	fmt.Println(canidate.name, " ", canidate.position)
	fmt.Println("  Votes ", canidate.numVotes)
}

// Returns the distortion of the election
func GetDistortion(winnterCost, optimalCost float64) float64 {
	return (winnterCost / optimalCost)
}

// Returns the social cost of a canidate
func getSocailCost(c Candidate, voters []Voter) float64 {
	cost := 0.0
	for i := 0; i < len(voters); i++ {
		cost += math.Abs(voters[i].position-c.position) * float64(voters[i].number)
	}
	return cost
}
