// File to simulate voting data and run the elections
package primary

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Generates a given number of "Empty canidates"
func initiateCanidates(numCandidates int, interval float64) []Candidate {
	var candidates []Candidate
	for i := 0; i < numCandidates; i++ {
		candidates = append(candidates, Candidate{Name: "Candidate" + string(rune(i)), Position: 0.0})
	}
	return candidates
}

// Places a given set of canidates onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeCandidates(candidates []Candidate, minPosition float64, maxPosition float64) []Candidate {
	for i := 0; i < len(candidates); i++ {

		candidates[i].Position = math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100
	}
	//Return the candidates with their new positions
	return candidates
}

// Generates a given number of "Empty voters"
func initiateVoters(numVoters int, interval float64) []Voter {
	var voters []Voter
	for i := 0; i < numVoters; i++ {
		voters = append(voters, Voter{Name: "Voter" + string(rune(i)), Number: 0.0, Position: 0.0})
	}
	return voters
}

// Places a given set of voters onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeVoters(numVoters int, minPosition float64, maxPosition float64, totalVoters int) []Voter {
	var voters []Voter
	//Create a set of 3 random numbers summing to 1
	numbers := generateRandomNumbers(numVoters, totalVoters)

	for i := 0; i < numVoters; i++ {
		voters = append(voters, Voter{Name: "Voter" + string(rune(i)), Number: numbers[i], Position: math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100})
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
		optimalCost, _ := DetermineOptimal(canidates, voters)
		canidates = Round(voters, canidates)
		//printCanidates(canidates)
		winnerCost := DetermineRoundWinner(canidates, totalVoters, voters)
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
