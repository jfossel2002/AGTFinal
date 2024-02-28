package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
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

// Generates a given number of "Empty canidates"
func initiateCanidates(numCandidates int, interval float64) []Candidate {
	var candidates []Candidate
	for i := 0; i < numCandidates; i++ {
		candidates = append(candidates, Candidate{name: "Candidate" + string(rune(i)), position: 0.0})
	}
	return candidates
}

// Places a given set of canidates onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeCandidates(candidates []Candidate, minPosition float64, maxPosition float64) []Candidate {
	for i := 0; i < len(candidates); i++ {

		candidates[i].position = math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100
	}
	//Return the candidates with their new positions
	return candidates
}

// Generates a given number of "Empty voters"
func initiateVoters(numVoters int, interval float64) []Voter {
	var voters []Voter
	for i := 0; i < numVoters; i++ {
		voters = append(voters, Voter{name: "Voter" + string(rune(i)), number: 0.0, position: 0.0})
	}
	return voters
}

// Places a given set of voters onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeVoters(numVoters int, minPosition float64, maxPosition float64, totalVoters int) []Voter {
	var voters []Voter
	//Create a set of 3 random numbers summing to 1
	numbers := generateRandomNumbers(numVoters, totalVoters)

	for i := 0; i < numVoters; i++ {
		voters = append(voters, Voter{name: "Voter" + string(rune(i)), number: numbers[i], position: math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100})
	}
	return voters
}

// generateRandomNumbers generates x random numbers that sum up to y.
func generateRandomNumbers(x int, y int) []int {
	rand.Seed(time.Now().UnixNano())

	var numbers []int
	sum := 0

	for i := 0; i < x-1; i++ {
		// Ensure the random number is less than the remaining difference to sum to `y`
		remaining := y - sum - (x - i - 1)
		if remaining <= 0 {
			numbers = append(numbers, 0)
		} else {
			n := rand.Intn(remaining/2) + 1
			sum += n
			numbers = append(numbers, n)
		}
	}

	// Add the final number to make the sum exactly `y`
	numbers = append(numbers, y-sum)

	return numbers
}

// Simulates a round by determining the number of votes each canidate gets
func round(voters []Voter, canidates []Candidate) []Candidate {
	//Reset the number of votes for each candidate
	for i := 0; i < len(canidates); i++ {
		canidates[i].numVotes = 0
	}
	//For each voter determine the closest candidate and save it
	for i := 0; i < len(voters); i++ {
		minDistance := 10000.0
		canidatePosition := 0
		for j := 0; j < len(canidates); j++ {
			distance := math.Abs(voters[i].position - canidates[j].position)
			if distance < minDistance {
				minDistance = distance
				canidatePosition = j
			}
		}
		//add that number of voters to the candidate
		canidates[canidatePosition].numVotes += voters[i].number

	}
	return canidates
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

func determineRoundWinner(candidates []Candidate, totalVoters int, voters []Voter) float64 {
	cost := -1.0
	//Base case
	for i := range candidates {
		if candidates[i].numVotes > totalVoters/2 {
			return getSocailCost(candidates[i], voters)
		}
	}

	//Recursive Case
	if cost == -1.0 {
		minVotes := totalVoters // Assuming totalVoters is higher than any candidate's numVotes can be
		candidatePosition := 0
		for j, candidate := range candidates {
			if candidate.numVotes < minVotes {
				minVotes = candidate.numVotes
				candidatePosition = j
			}
		}

		// Remove the candidate with the least votes
		candidates = append(candidates[:candidatePosition], candidates[candidatePosition+1:]...)

		// Redetermine votes based on voters' next preferences
		//printCanidates(candidates)
		candidates = round(voters, candidates)

		// Recursive call to handle the next round
		return determineRoundWinner(candidates, totalVoters, voters)
	}

	return cost
}

// Finds the optimal canidate based on the social cost
func determineOptimal(Candidates []Candidate, voters []Voter) (float64, int) {
	minCost := 1000000.0
	canidatePosition := 0
	for j := 0; j < len(Candidates); j++ {
		getSocailCost(Candidates[j], voters)
		distance := getSocailCost(Candidates[j], voters)
		if distance < minCost {
			minCost = distance
			canidatePosition = j
		}
	}
	//fmt.Println("The optimal canidate is ", Candidates[canidatePosition].name, " ", Candidates[canidatePosition].position, " with a cost of ", minCost)
	return minCost, canidatePosition

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

// Runs a given number of scenarios
func runScenario(numRuns int, numCandidates int, maxPosition float64, minPosition float64, totalVoters int) {
	maxDistortion := 0.0
	maxDistortionCanidates := []Candidate{}
	maxDistortionVoters := []Voter{}
	for i := 0; i < numRuns; i++ {
		canidates := initiateCanidates(numCandidates, maxPosition)
		canidates = distributeCandidates(canidates, minPosition, maxPosition)
		//fmt.Println(canidates)
		voters := initiateVoters(5, maxPosition)
		voters = distributeVoters(5, minPosition, maxPosition, totalVoters)
		//printVoters(voters)
		optimalCost, _ := determineOptimal(canidates, voters)
		canidates = round(voters, canidates)
		//printCanidates(canidates)
		winnerCost := determineRoundWinner(canidates, totalVoters, voters)
		//fmt.Println("The winner cost is ", winnerCost)
		//fmt.Println("The optimal cost is ", optimalCost)
		distortion := GetDistortion(winnerCost, optimalCost)
		if distortion > maxDistortion {
			maxDistortion = distortion
			maxDistortionCanidates = canidates
			maxDistortionVoters = voters
		}
		if distortion > 3 {
			fmt.Println("The distortion is ", distortion)
			printCanidates(canidates)
			printVoters(voters)
		}
	}
	fmt.Println("The max distortion is ", maxDistortion)
	printCanidates(maxDistortionCanidates)
	printVoters(maxDistortionVoters)
}

func runSpecificInstance() {
	totalVoters := 7
	/*canidates := []Candidate{
		{"A", 0.11, 0},
		{"B", 0.11, 0},
		{"C", 3.78, 0},
		{"D", 5.63, 0},
		{"E", 9.92, 0},
	}

	voters := []Voter{
		{"V1", 23, 1.51},
		{"V2", 23, 4.45},
		{"V3", 4, 5.45},
	}*/

	canidates := []Candidate{
		{"W", 0, 0},
		{"X", 4, 0},
		{"C", 8, 0},
	}

	voters := []Voter{
		{"X=2", 1, 2},
		{"X=3", 2, 3},
		{"X=4", 2, 4},
		{"X=6", 2, 6},
	}
	optimalCost, _ := determineOptimal(canidates, voters)
	canidates = round(voters, canidates)
	//printCanidates(canidates)
	winnerCost := determineRoundWinner(canidates, totalVoters, voters)
	//printCanidates(canidates)
	fmt.Println("The winner cost is ", winnerCost)
	fmt.Println("The optimal cost is ", optimalCost)
	distortion := GetDistortion(winnerCost, optimalCost)
	fmt.Println("The distortion is ", distortion)
	//printCanidates(canidates)
	printVoters(voters)

}

func main() {
	minPosition := 0.0
	maxPosition := 1.0
	numCandidates := 4
	totalVoters := 50
	//runSpecificInstance()
	runScenario(100000, numCandidates, maxPosition, minPosition, totalVoters)
}
