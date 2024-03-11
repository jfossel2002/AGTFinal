// File to run a specific scenario
package primary

import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"
)

func RunSpecificInstance(canidates []voting_systems.Candidate, voters []voting_systems.Voter) {
	totalVoters := 50

	//Print the cost for each canidate
	voting_systems.PrintAllCosts(canidates, voters)
	optimalCost, opt_canidate := voting_systems.DetermineOptimalCanidate(canidates, voters)
	fmt.Println("\nThe optimal canidate is ", opt_canidate.Name, " ", opt_canidate.Position)
	fmt.Println("The optimal cost is ", optimalCost)
	winner := voting_systems.InitiateSTV(canidates, totalVoters, voters)
	fmt.Println("The winner From STV is ", winner.Name, " ", winner.Position)
	winnerCost := voting_systems.GetSocailCost(winner, voters)
	fmt.Println("The winner From STV cost is ", winnerCost)

	//Add Borda count determination here
	bordaWinner := voting_systems.CalculateBordaWinner(canidates, voters)
	fmt.Println("The winner From Borda Count is:", bordaWinner.Name, " ", bordaWinner.Position)
	// Copeland count determination here
	copelandWinner := voting_systems.DetermineCopelandWinner(canidates, voters)
	fmt.Println("The winner From Copeland is:", copelandWinner.Name, " ", copelandWinner.Position)

	//Plurality count determination here
	pluralityWinner := voting_systems.InitiatePlurality(canidates, voters)
	fmt.Println("The winner From Plurality is:", pluralityWinner.Name, " ", pluralityWinner.Position)

	//Plurality Veto determination here
	pluralityVetoWinner := voting_systems.InitiatePluralityVeto(canidates, voters)
	fmt.Println("The winner From Plurality Veto is:", pluralityVetoWinner.Name, " ", pluralityVetoWinner.Position)

}
