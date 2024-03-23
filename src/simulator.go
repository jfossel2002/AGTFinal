// File to simulate voting data and run the elections
package primary

import (
	voting_systems "AGT_Midterm/src/systems"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
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
			STVWinner, _, _ := voting_systems.InitiateSTV(canidates, voters)
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
			bordaWinner, _ := voting_systems.CalculateBordaWinner(canidates, voters)
			bordaWinnerCost := voting_systems.GetSocailCost(bordaWinner, voters)
			bordaDistortion := voting_systems.GetDistortion(bordaWinnerCost, optimalCost)
			if bordaDistortion > maxBordaDistortion {
				maxBordaDistortion = bordaDistortion
				maxBordaDistortionCanidates = canidates
				maxBordaDistortionVoters = voters
			}

		}
		if votingSystem == "Plurality" || votingSystem == "All" {
			pluralityWinner, _ := voting_systems.InitiatePlurality(canidates, voters)
			pluralityWinnerCost := voting_systems.GetSocailCost(pluralityWinner, voters)
			pluralityDistortion := voting_systems.GetDistortion(pluralityWinnerCost, optimalCost)
			if pluralityDistortion > maxPluralityDistortion {
				maxPluralityDistortion = pluralityDistortion
				maxPluralityDistortionCanidates = canidates
				maxPluralityDistortionVoters = voters
			}
		}
		if votingSystem == "Copeland" || votingSystem == "All" {
			copelandWinner, _ := voting_systems.DetermineCopelandWinner(canidates, voters)
			copelandWinnerCost := voting_systems.GetSocailCost(copelandWinner, voters)
			copelandDistortion := voting_systems.GetDistortion(copelandWinnerCost, optimalCost)
			if copelandDistortion > maxCopelandDistortion {
				maxCopelandDistortion = copelandDistortion
				maxCopelandDistortionCanidates = canidates
				maxCopelandDistortionVoters = voters
			}
		}
		if votingSystem == "Plurality Veto" || votingSystem == "All" {
			pluralityVetoWinner, _, _ := voting_systems.InitiatePluralityVeto(canidates, voters)
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
		winnerSTVCan, _, _ := voting_systems.InitiateSTV(maxSTVDistortionCanidates, maxSTVDistortionVoters)
		fmt.Println("The winner is ", winnerSTVCan)
		voting_systems.PrintCanidates(maxSTVDistortionCanidates)
		voting_systems.PrintVoters(maxSTVDistortionVoters)
		saveCanidates(maxSTVDistortionCanidates, "STV-Canidates-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		saveVoters(maxSTVDistortionVoters, "STV-Voters-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		return "Max STV Distortion: " + fmt.Sprintf("%f", maxSTVDistortion)
	}
	if votingSystem == "Borda Count" || votingSystem == "All" {
		fmt.Println("\nThe max Borda Count distortion is ", maxBordaDistortion)
		_, optBordaCan := voting_systems.DetermineOptimalCanidate(maxBordaDistortionCanidates, maxBordaDistortionVoters)
		winnerBordaCan, _ := voting_systems.CalculateBordaWinner(maxBordaDistortionCanidates, maxBordaDistortionVoters)
		fmt.Println("The optimal canidate is ", optBordaCan)
		fmt.Println("The winner is ", winnerBordaCan)
		voting_systems.PrintCanidates(maxBordaDistortionCanidates)
		voting_systems.PrintVoters(maxBordaDistortionVoters)
		saveCanidates(maxBordaDistortionCanidates, "Borda-Canidates-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		saveVoters(maxBordaDistortionVoters, "Borda-Voters-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")

		return "Max Borda Distortion: " + fmt.Sprintf("%f", maxBordaDistortion)
	}
	if votingSystem == "Plurality" || votingSystem == "All" {
		fmt.Println("\nThe max Plurality distortion is ", maxPluralityDistortion)
		_, optPluralityCan := voting_systems.DetermineOptimalCanidate(maxPluralityDistortionCanidates, maxPluralityDistortionVoters)
		winnerPluralityCan, _ := voting_systems.InitiatePlurality(maxPluralityDistortionCanidates, maxPluralityDistortionVoters)
		fmt.Println("The optimal canidate is ", optPluralityCan)
		fmt.Println("The winner is ", winnerPluralityCan)
		voting_systems.PrintCanidates(maxPluralityDistortionCanidates)
		voting_systems.PrintVoters(maxPluralityDistortionVoters)
		saveCanidates(maxPluralityDistortionCanidates, "Plurality-Canidates-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		saveVoters(maxPluralityDistortionVoters, "Plurality-Voters-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		return "Max Plurality Distortion: " + fmt.Sprintf("%f", maxPluralityDistortion)
	}
	if votingSystem == "Copeland" || votingSystem == "All" {
		fmt.Println("\nThe max Copeland distortion is ", maxCopelandDistortion)
		_, optCopelandCan := voting_systems.DetermineOptimalCanidate(maxCopelandDistortionCanidates, maxCopelandDistortionVoters)
		winnerCopelandCan, _ := voting_systems.DetermineCopelandWinner(maxCopelandDistortionCanidates, maxCopelandDistortionVoters)
		fmt.Println("The optimal canidate is ", optCopelandCan)
		fmt.Println("The winner is ", winnerCopelandCan)
		voting_systems.PrintCanidates(maxCopelandDistortionCanidates)
		voting_systems.PrintVoters(maxCopelandDistortionVoters)
		saveCanidates(maxCopelandDistortionCanidates, "Copeland-Canidates-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		saveVoters(maxCopelandDistortionVoters, "Copeland-Voters-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		return "Max Copeland Distortion: " + fmt.Sprintf("%f", maxCopelandDistortion)
	}
	if votingSystem == "Plurality Veto" || votingSystem == "All" {
		fmt.Println("\nThe max Plurality Veto distortion is ", maxPluralityVetoDistortion)
		_, optPluralityVetoCan := voting_systems.DetermineOptimalCanidate(maxPluralityVetoDistortionCanidates, maxPluralityVetoDistortionVoters)
		winnerPluralityVetoCan, _, _ := voting_systems.InitiatePluralityVeto(maxPluralityVetoDistortionCanidates, maxPluralityVetoDistortionVoters)
		fmt.Println("The optimal canidate is ", optPluralityVetoCan)
		fmt.Println("The winner is ", winnerPluralityVetoCan)
		voting_systems.PrintCanidates(maxPluralityVetoDistortionCanidates)
		voting_systems.PrintVoters(maxPluralityVetoDistortionVoters)
		saveCanidates(maxPluralityVetoDistortionCanidates, "PVeto-Canidates-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		saveVoters(maxPluralityVetoDistortionVoters, "PVeto-Voters-"+fmt.Sprintf("%.2f", maxSTVDistortion)+"-"+time.Now().Format("2006-01-02-15:04:05")+".json")
		return "Max PLurality Veto Distortion: " + fmt.Sprintf("%f", maxPluralityVetoDistortion)
	}
	return "Done"
}

func saveCanidates(candidates []voting_systems.Candidate, filename string) {
	// Path where the JSON files will be saved
	basePath := "Jsons/Candidates"

	// Ensure the directory structure exists
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		panic(err) // Handle the error properly in a real scenario
	}

	// Sanitize the Name field of each candidate
	for i := range candidates {
		candidates[i].Name = SanitizeString(candidates[i].Name)
	}

	// Combine the base path with the filename to get the full path
	fullPath := filepath.Join(basePath, filename)

	// Marshal the candidates slice to JSON
	file, err := json.MarshalIndent(candidates, "", " ")
	if err != nil {
		panic(err) // Handle the error properly in a real scenario
	}

	// Write the file to the specified path
	err = ioutil.WriteFile(fullPath, file, 0644)
	if err != nil {
		panic(err) // Handle the error properly in a real scenario
	}
}

// SanitizeString removes ASCII control characters from a string
func SanitizeString(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1 // Drop control characters
		}
		return r
	}, s)
}

func saveVoters(voters []voting_systems.Voter, filename string) {
	// Path where the JSON files will be saved
	basePath := "Jsons/Voters"

	// Ensure the directory structure exists
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		// Handle the error properly in a real scenario
		panic(err)
	}

	// Sanitize the Name field of each voter
	for i := range voters {
		voters[i].Name = SanitizeString(voters[i].Name)
	}

	// Combine the base path with the filename to get the full path
	fullPath := filepath.Join(basePath, filename)

	// Marshal the voters slice to JSON
	file, err := json.MarshalIndent(voters, "", " ")
	if err != nil {
		// Handle the error properly in a real scenario
		panic(err)
	}

	// Write the file to the specified path
	err = ioutil.WriteFile(fullPath, file, 0644)
	if err != nil {
		// Handle the error properly in a real scenario
		panic(err)
	}
}
