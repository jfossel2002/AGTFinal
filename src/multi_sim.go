package primary

type MultiInput struct {
	NumCandidates int
	NumVoters     int
}

type MultiOutput struct {
	Result            string
	CandidateFileName string
	VoterFileName     string
	MaxDistortion     float64
	AverageDistortion float64
}

func Multi_sim(minCandidates int, maxCandidates int, minVoters int, maxVoters int, numRuns int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) map[MultiInput]MultiOutput {
	//Create a map to store the results
	results := make(map[MultiInput]MultiOutput)
	for candidates := minCandidates; candidates <= maxCandidates; candidates++ {
		for voters := minVoters; voters <= maxVoters; voters++ {
			result, candidateFileName, voterFileName, maxDistortion, averageDistortion := RunScenario(numRuns, candidates, voters, maxPosition, minPosition, totalVoters, votingSystem)
			output := MultiOutput{
				Result:            result,
				CandidateFileName: candidateFileName,
				VoterFileName:     voterFileName,
				MaxDistortion:     maxDistortion,
				AverageDistortion: averageDistortion,
			}
			results[MultiInput{candidates, voters}] = output
		}
	}
	return results
}
