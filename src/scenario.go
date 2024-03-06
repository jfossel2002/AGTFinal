// File to run a specific scenario
package primary

import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"
)

func RunSpecificInstance() {
	totalVoters := 50
	canidates := []voting_systems.Candidate{
		{Name: "A", Position: 0.11, NumVotes: 0},
		{Name: "B", Position: 0.11, NumVotes: 0},
		{Name: "C", Position: 3.78, NumVotes: 0},
		{Name: "D", Position: 5.63, NumVotes: 0},
		{Name: "E", Position: 9.92, NumVotes: 0},
	}

	voters := []voting_systems.Voter{
		{Name: "V1", Number: 23, Position: 1.51},
		{Name: "V2", Number: 23, Position: 4.45},
		{Name: "V3", Number: 4, Position: 5.45},
	}

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

}
