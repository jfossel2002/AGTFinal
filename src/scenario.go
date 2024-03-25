package primary

/*
* This file contains the functions for running a speicifc scenario (a specific set of candidates and voters for a voting system)
 */
import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"
)

// RunSpecificInstance runs a specific scenario
// Takes in a slice of candidates and a slice of voters
// Prints the optimal candidate, optimal cost, winner from STV, winner from Borda Count, winner from Copeland, winner from Plurality, and winner from Plurality Veto
func RunSpecificInstance(canidates []voting_systems.Candidate, voters []voting_systems.Voter) {
	//Print the cost for each canidate
	voting_systems.PrintAllCosts(canidates, voters)
	optimalCost, opt_canidate := voting_systems.DetermineOptimalCanidate(canidates, voters)
	fmt.Println("\nThe optimal canidate is ", opt_canidate.Name, " ", opt_canidate.Position)
	fmt.Println("The optimal cost is ", optimalCost)
	winner, _, _ := voting_systems.InitiateSTV(canidates, voters)
	fmt.Println("The winner From STV is ", winner.Name, " ", winner.Position)
	winnerCost := voting_systems.GetSocailCost(winner, voters)
	fmt.Println("The winner From STV cost is ", winnerCost)

	//Add Borda count determination here
	bordaWinner, _ := voting_systems.CalculateBordaWinner(canidates, voters)
	fmt.Println("The winner From Borda Count is:", bordaWinner.Name, " ", bordaWinner.Position)
	// Copeland count determination here
	copelandWinner, _ := voting_systems.DetermineCopelandWinner(canidates, voters)
	fmt.Println("The winner From Copeland is:", copelandWinner.Name, " ", copelandWinner.Position)

	//Plurality count determination here
	pluralityWinner, _ := voting_systems.InitiatePlurality(canidates, voters)
	fmt.Println("The winner From Plurality is:", pluralityWinner.Name, " ", pluralityWinner.Position)

	//Plurality Veto determination here
	pluralityVetoWinner, _, _ := voting_systems.InitiatePluralityVeto(canidates, voters)
	fmt.Println("The winner From Plurality Veto is:", pluralityVetoWinner.Name, " ", pluralityVetoWinner.Position)

}
