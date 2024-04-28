/*
* This file contains the functions to read data from a JSON file, and the functions to calculate the social cost of a candidate, the distortion of the election, and the optimal candidate.
* Contains helper functions and sutrcutres that are shared accross all voting methods
 */
package voting_systems

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

// Define the Candidate struct
type Candidate struct {
	Name     string  `json:"Name"`
	Position float64 `json:"Position"`
	NumVotes int     `json:"NumVotes"`
}

// Define the Voter struct
type Voter struct {
	Name     string  `json:"Name"`
	Number   int     `json:"Number"`
	Position float64 `json:"Position"`
}

// Returns the distortion of the election
func GetDistortion(winnterCost, optimalCost float64) float64 {
	if optimalCost == 0 {
		return 1
	}
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
	return minCost, Candidates[canidatePosition]
}

// Prints the cost of all canidates
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

// Function to read data from a JSON file, user specifies the data type (Candidate or Voter)
// Returns a slice of Candidates or Voters based on the dataType
func ReadFromFile(filePath string, dataType string) (interface{}, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if dataType == "Candidate" {
		// If the data type is Candidate, unmarshal into []Candidate
		var candidates []Candidate
		err = json.Unmarshal(byteValue, &candidates)
		if err != nil {
			return nil, err
		}
		return candidates, nil
	} else if dataType == "Voter" {
		// If the data type is Voter, unmarshal into []Voter
		var voters []Voter
		err = json.Unmarshal(byteValue, &voters)
		if err != nil {
			return nil, err
		}
		return voters, nil
	}

	// If an unsupported dataType is provided
	return nil, fmt.Errorf("unsupported data type")
}
