// File to handle main functions
package main

import (
	primary "AGT_Midterm/src"
	voting_systems "AGT_Midterm/src/systems"
)

func main() {
	canidatesFilePath := "canidates.json"
	votersFilePath := "voters.json"

	candidatesInterface, err := voting_systems.ReadFromFile(canidatesFilePath, "Candidate")
	candidates := candidatesInterface.([]voting_systems.Candidate)

	votersInterface, err := voting_systems.ReadFromFile(votersFilePath, "Voter")
	voters := votersInterface.([]voting_systems.Voter)

	if err != nil {
		panic(err)
	}

	primary.RunSpecificInstance(candidates, voters)
	/*minPosition := 0.0
	maxPosition := 1.0
	numCandidates := 4
	totalVoters := 50

	primary.RunScenario(1000000, numCandidates, maxPosition, minPosition, totalVoters)*/

}
