package primary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"
)

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

type MultiResults struct {
	Results map[string]MultiOutput
}

func Multi_sim(minCandidates int, maxCandidates int, minVoters int, maxVoters int, numRuns int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) (map[MultiInput]MultiOutput, string) {
	results := make(map[MultiInput]MultiOutput)
	for candidates := minCandidates; candidates <= maxCandidates; candidates++ {
		for voters := minVoters; voters <= maxVoters; voters++ {
			result, candidateFileName, voterFileName, maxDistortion, averageDistortion := RunScenario(numRuns, candidates, voters, maxPosition, minPosition, totalVoters, votingSystem)

			if math.IsNaN(maxDistortion) {
				maxDistortion = 0
			}
			if math.IsNaN(averageDistortion) {
				averageDistortion = 0
			}

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
	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")
	fileName := "multi_sim_" + currentDate + "_" + currentTime + ".json"
	if err := saveJSON(convertToMultiResults(results, fileName), fileName); err != nil {
		fmt.Println("Error saving JSON file:", err)
	}
	return results, fileName
}

func convertToMultiResults(results map[MultiInput]MultiOutput, fileName string) MultiResults {
	multiResults := MultiResults{
		Results: make(map[string]MultiOutput),
	}
	for input, output := range results {
		key := fmt.Sprintf("%d_%d", input.NumCandidates, input.NumVoters)
		multiResults.Results[key] = output
	}
	return multiResults
}

func saveJSON(results MultiResults, fileName string) error {
	file, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadAndGraphMultiResults(fileName string) {
	results := MultiResults{}
	err := readJSON(fileName, &results)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	for key, value := range results.Results {
		fmt.Println("Key:", key)
		fmt.Println("Value:", value)
	}
	//Key: 4_5
	//Value: {Max STV Distortion: 1.237255 STV-Canidates-1.24-2024-04-28-16:04:25.json STV-Voters-1.24-2024-04-28-16:04:25.json 1.2372549019607841 1.054829236488238}
	GraphResults()
}

func readJSON(fileName string, results *MultiResults) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, results)
	if err != nil {
		return err
	}
	return nil
}

func GraphResults() {

}
