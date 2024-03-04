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

	optimalCost, opt_canidate := voting_systems.DetermineOptimalCanidate(canidates, voters)
	fmt.Println("The optimal canidate is ", opt_canidate.Name, " ", opt_canidate.Position)
	fmt.Println("The optimal cost is ", optimalCost)
	winner := voting_systems.InitiateSTV(canidates, totalVoters, voters)
	fmt.Println("The winner is ", winner.Name, " ", winner.Position)
	winnerCost := voting_systems.GetSocailCost(winner, voters)
	fmt.Println("The winner cost is ", winnerCost)

}
