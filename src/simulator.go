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
func RunScenario(numRuns int, numCandidates int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) string {
	maxSTVDistortion := -1.0
	maxBordaDistortion := -1.0
	maxSTVDistortionCanidates := []voting_systems.Candidate{}
	maxSTVDistortionVoters := []voting_systems.Voter{}
	maxBordaDistortionCanidates := []voting_systems.Candidate{}
	maxBordaDistortionVoters := []voting_systems.Voter{}
	maxPluralityDistortion := -1.0
	maxPluralityDistortionCanidates := []voting_systems.Candidate{}
	maxPluralityDistortionVoters := []voting_systems.Voter{}
	maxCopelandDistortion := -1.0
	maxCopelandDistortionCanidates := []voting_systems.Candidate{}
	maxCopelandDistortionVoters := []voting_systems.Voter{}
	maxPluralityVetoDistortion := -1.0
	maxPluralityVetoDistortionCanidates := []voting_systems.Candidate{}
	maxPluralityVetoDistortionVoters := []voting_systems.Voter{}
	for i := 0; i < numRuns; i++ {
		canidates := initiateCanidates(numCandidates, maxPosition)
		canidates = distributeCandidates(canidates, minPosition, maxPosition)
		//fmt.Println(canidates)
		voters := initiateVoters(5, maxPosition)
		voters = distributeVoters(5, minPosition, maxPosition, totalVoters)
		//fmt.Println(voters)
		optimalCost, _ := voting_systems.DetermineOptimalCanidate(canidates, voters)
		if votingSystem == "STV" || votingSystem == "All" {
			STVWinner := voting_systems.InitiateSTV(canidates, totalVoters, voters)
			STVWinnerCost := voting_systems.GetSocailCost(STVWinner, voters)
			STVDistortion := voting_systems.GetDistortion(STVWinnerCost, optimalCost)
			if STVDistortion > maxSTVDistortion {
				maxSTVDistortion = STVDistortion
				maxSTVDistortionCanidates = canidates
				maxSTVDistortionVoters = voters
			}
			if STVDistortion > 3 {
				fmt.Println("The STV distortion is ", STVDistortion)
				voting_systems.PrintCanidates(canidates)
				voting_systems.PrintVoters(voters)
			}

		}
		if votingSystem == "Borda Count" || votingSystem == "All" {
			bordaWinner := voting_systems.CalculateBordaWinner(canidates, voters)
			bordaWinnerCost := voting_systems.GetSocailCost(bordaWinner, voters)
			bordaDistortion := voting_systems.GetDistortion(bordaWinnerCost, optimalCost)
			if bordaDistortion > maxBordaDistortion {
				maxBordaDistortion = bordaDistortion
				maxBordaDistortionCanidates = canidates
				maxBordaDistortionVoters = voters
			}

		}
		if votingSystem == "Plurality" || votingSystem == "All" {
			pluralityWinner := voting_systems.InitiatePlurality(canidates, voters)
			pluralityWinnerCost := voting_systems.GetSocailCost(pluralityWinner, voters)
			pluralityDistortion := voting_systems.GetDistortion(pluralityWinnerCost, optimalCost)
			if pluralityDistortion > maxPluralityDistortion {
				maxPluralityDistortion = pluralityDistortion
				maxPluralityDistortionCanidates = canidates
				maxPluralityDistortionVoters = voters
			}
		}
		if votingSystem == "Copeland" || votingSystem == "All" {
			copelandWinner := voting_systems.DetermineCopelandWinner(canidates, voters)
			copelandWinnerCost := voting_systems.GetSocailCost(copelandWinner, voters)
			copelandDistortion := voting_systems.GetDistortion(copelandWinnerCost, optimalCost)
			if copelandDistortion > maxCopelandDistortion {
				maxCopelandDistortion = copelandDistortion
				maxCopelandDistortionCanidates = canidates
				maxCopelandDistortionVoters = voters
			}
		}
		if votingSystem == "Plurality Veto" || votingSystem == "All" {
			pluralityVetoWinner := voting_systems.InitiatePluralityVeto(canidates, voters)
			pluralityVetoWinnerCost := voting_systems.GetSocailCost(pluralityVetoWinner, voters)
			pluralityVetoDistortion := voting_systems.GetDistortion(pluralityVetoWinnerCost, optimalCost)
			if pluralityVetoDistortion > maxPluralityVetoDistortion {
				maxPluralityVetoDistortion = pluralityVetoDistortion
				maxPluralityVetoDistortionCanidates = canidates
				maxPluralityVetoDistortionVoters = voters
			}
		}

		//fmt.Println("The winner cost is ", winnerCost)
		//fmt.Println("The optimal cost is ", optimalCost)

	}
	if votingSystem == "STV" || votingSystem == "All" {
		fmt.Println("The max STV distortion is ", maxSTVDistortion)
		_, optSTVCan := voting_systems.DetermineOptimalCanidate(maxSTVDistortionCanidates, maxSTVDistortionVoters)
		fmt.Println("The optimal STV canidate is ", optSTVCan)
		winnerSTVCan := voting_systems.InitiateSTV(maxSTVDistortionCanidates, totalVoters, maxSTVDistortionVoters)
		fmt.Println("The winner is ", winnerSTVCan)
		voting_systems.PrintCanidates(maxSTVDistortionCanidates)
		voting_systems.PrintVoters(maxSTVDistortionVoters)
		return "Max STV Distortion: " + fmt.Sprintf("%f", maxSTVDistortion)
	}
	if votingSystem == "Borda Count" || votingSystem == "All" {
		fmt.Println("\nThe max Borda Count distortion is ", maxBordaDistortion)
		_, optBordaCan := voting_systems.DetermineOptimalCanidate(maxBordaDistortionCanidates, maxBordaDistortionVoters)
		winnerBordaCan := voting_systems.CalculateBordaWinner(maxBordaDistortionCanidates, maxBordaDistortionVoters)
		fmt.Println("The optimal canidate is ", optBordaCan)
		fmt.Println("The winner is ", winnerBordaCan)
		voting_systems.PrintCanidates(maxBordaDistortionCanidates)
		voting_systems.PrintVoters(maxBordaDistortionVoters)
		return "Max Borda Distortion: " + fmt.Sprintf("%f", maxBordaDistortion)
	}
	if votingSystem == "Plurality" || votingSystem == "All" {
		fmt.Println("\nThe max Plurality distortion is ", maxPluralityDistortion)
		_, optPluralityCan := voting_systems.DetermineOptimalCanidate(maxPluralityDistortionCanidates, maxPluralityDistortionVoters)
		winnerPluralityCan := voting_systems.InitiatePlurality(maxPluralityDistortionCanidates, maxPluralityDistortionVoters)
		fmt.Println("The optimal canidate is ", optPluralityCan)
		fmt.Println("The winner is ", winnerPluralityCan)
		voting_systems.PrintCanidates(maxPluralityDistortionCanidates)
		voting_systems.PrintVoters(maxPluralityDistortionVoters)
		return "Max Plurality Distortion: " + fmt.Sprintf("%f", maxPluralityDistortion)
	}
	if votingSystem == "Copeland" || votingSystem == "All" {
		fmt.Println("\nThe max Copeland distortion is ", maxCopelandDistortion)
		_, optCopelandCan := voting_systems.DetermineOptimalCanidate(maxCopelandDistortionCanidates, maxCopelandDistortionVoters)
		winnerCopelandCan := voting_systems.DetermineCopelandWinner(maxCopelandDistortionCanidates, maxCopelandDistortionVoters)
		fmt.Println("The optimal canidate is ", optCopelandCan)
		fmt.Println("The winner is ", winnerCopelandCan)
		voting_systems.PrintCanidates(maxCopelandDistortionCanidates)
		voting_systems.PrintVoters(maxCopelandDistortionVoters)
		return "Max Copeland Distortion: " + fmt.Sprintf("%f", maxCopelandDistortion)
	}
	if votingSystem == "Plurality Veto" || votingSystem == "All" {
		fmt.Println("\nThe max Plurality Veto distortion is ", maxPluralityVetoDistortion)
		_, optPluralityVetoCan := voting_systems.DetermineOptimalCanidate(maxPluralityVetoDistortionCanidates, maxPluralityVetoDistortionVoters)
		winnerPluralityVetoCan := voting_systems.InitiatePluralityVeto(maxPluralityVetoDistortionCanidates, maxPluralityVetoDistortionVoters)
		fmt.Println("The optimal canidate is ", optPluralityVetoCan)
		fmt.Println("The winner is ", winnerPluralityVetoCan)
		voting_systems.PrintCanidates(maxPluralityVetoDistortionCanidates)
		voting_systems.PrintVoters(maxPluralityVetoDistortionVoters)
		return "Max PLurality Veto Distortion: " + fmt.Sprintf("%f", maxPluralityVetoDistortion)
	}
	return "Done"
}
