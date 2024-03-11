// File to simulate voting data and run the elections
package primary

import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func initiateCanidates(numCandidates int, interval float64) []voting_systems.Candidate {
	var candidates []voting_systems.Candidate
	for i := 0; i < numCandidates; i++ {
		candidates = append(candidates, voting_systems.Candidate{Name: "Candidate" + string(rune(i)), Position: 0.0})
	}
	return candidates
}

// Places a given set of canidates onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeCandidates(candidates []voting_systems.Candidate, minPosition float64, maxPosition float64) []voting_systems.Candidate {
	for i := 0; i < len(candidates); i++ {

		candidates[i].Position = math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100
	}
	//Return the candidates with their new positions
	return candidates
}

// Generates a given number of "Empty voters"
func initiateVoters(numVoters int, interval float64) []voting_systems.Voter {
	var voters []voting_systems.Voter
	for i := 0; i < numVoters; i++ {
		voters = append(voters, voting_systems.Voter{Name: "Voter" + string(rune(i)), Number: 0.0, Position: 0.0})
	}
	return voters
}

// Places a given set of voters onto the interval 0-1 randomly with up to 2 decimal places accuracy
func distributeVoters(numVoters int, minPosition float64, maxPosition float64, totalVoters int) []voting_systems.Voter {
	var voters []voting_systems.Voter
	//Create a set of 3 random numbers summing to 1
	numbers := generateRandomNumbers(numVoters, totalVoters)

	for i := 0; i < numVoters; i++ {
		voters = append(voters, voting_systems.Voter{Name: "Voter" + string(rune(i)), Number: numbers[i], Position: math.Round((minPosition+(maxPosition-minPosition)*rand.Float64())*100) / 100})
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
func RunScenario(numRuns int, numCandidates int, maxPosition float64, minPosition float64, totalVoters int) {
	maxDistortion := 0.0
	maxDistortionCanidates := []voting_systems.Candidate{}
	maxDistortionVoters := []voting_systems.Voter{}
	for i := 0; i < numRuns; i++ {
		canidates := initiateCanidates(numCandidates, maxPosition)
		canidates = distributeCandidates(canidates, minPosition, maxPosition)
		//fmt.Println(canidates)
		voters := initiateVoters(5, maxPosition)
		voters = distributeVoters(5, minPosition, maxPosition, totalVoters)
		//printVoters(voters)
		optimalCost, _ := voting_systems.DetermineOptimalCanidate(canidates, voters)
		STVWinner := voting_systems.InitiateSTV(canidates, totalVoters, voters)
		winnerCost := voting_systems.GetSocailCost(STVWinner, voters)
		//fmt.Println("The winner cost is ", winnerCost)
		//fmt.Println("The optimal cost is ", optimalCost)
		distortion := voting_systems.GetDistortion(winnerCost, optimalCost)
		if distortion > maxDistortion {
			maxDistortion = distortion
			maxDistortionCanidates = canidates
			maxDistortionVoters = voters
		}
		if distortion > 3 {
			fmt.Println("The distortion is ", distortion)
			voting_systems.PrintCanidates(canidates)
			voting_systems.PrintVoters(voters)
		}
	}
	fmt.Println("The max distortion is ", maxDistortion)
	voting_systems.PrintCanidates(maxDistortionCanidates)
	voting_systems.PrintVoters(maxDistortionVoters)
}
