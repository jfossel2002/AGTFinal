package primary

/*
This file contains the functions to run a simulation over a range as well as the graph generation logic for the multi-simulation.
*/
import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Structs for the multi-simulation
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

// Function to run a multi-simulation over a range of candidates and voters
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

// Function to convert the map of results to a MultiResults struct
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

// Function to save the results of the multi-simulation to a JSON file
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

// Function to read the results of a multi-simulation from a JSON file then graph the results
func ReadAndGraphMultiResults(fileName string, isMax bool, isCandidates bool) string {
	results := MultiResults{}
	err := readJSON(fileName, &results)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return "ERROR"
	}

	allPairs := MultiParse(results, true)
	return MultiGraph(allPairs, results, isMax, isCandidates)
}

// Function to generate a random color
func RandomColor() color.RGBA {
	rand.Seed(time.Now().UnixNano())
	return color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255, // Fully opaque
	}
}

// Function to generate a multi-graph i,e, a graph with multiple lines
func MultiGraph(mapData []int, results MultiResults, isMax bool, isCandidates bool) string {
	p := plot.New()

	if isMax {
		if isCandidates {
			p.Title.Text = "Max Distortion on Candidates"
		} else {
			p.Title.Text = "Max Distortion on Voters"
		}
	} else {
		if isCandidates {
			p.Title.Text = "Average Distortion on Candidates"
		} else {
			p.Title.Text = "Average Distortion on Voters"
		}
	}
	if isCandidates {
		p.X.Label.Text = "Voters"
	} else {
		p.X.Label.Text = "Candidates"
	}
	p.Y.Label.Text = "Distortion"

	for value := range mapData {
		pts := make(plotter.XYs, 0)
		parsedData := make(map[int]float64)

		if isCandidates {
			if isMax {
				parsedData = ParseCandidatesOrVoters(results, value+1, true, true)
			} else {
				parsedData = ParseCandidatesOrVoters(results, value+1, true, false)
			}
		} else {
			if isMax {
				parsedData = ParseCandidatesOrVoters(results, value+1, false, true)
			} else {
				parsedData = ParseCandidatesOrVoters(results, value+1, false, false)
			}
		}

		keys := make([]int, 0, len(parsedData))
		for key := range parsedData {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		for _, key := range keys {
			val := parsedData[key]
			pts = append(pts, plotter.XY{X: float64(key), Y: val})
		}

		lpLine, lpPoints, err := plotter.NewLinePoints(pts)
		if err != nil {
			panic(err)
		}
		color := RandomColor()
		lpLine.Color = color
		lpPoints.Color = color

		p.Add(lpLine, lpPoints)
		if isCandidates {
			p.Legend.Add(fmt.Sprintf("Candidates %d", value+1), lpLine, lpPoints)
		} else {
			p.Legend.Add(fmt.Sprintf("Voters %d", value+1), lpLine, lpPoints)
		}
	}

	p.Legend.Top = false

	// Save the plot to a PNG file.
	CurrentData := time.Now().Format("2006-01-02")
	CurrentTime := time.Now().Format("15:04:05")
	fileName := ""
	if isMax {
		if isCandidates {
			fileName = "MaxDistortionCandidates_" + CurrentData + "_" + CurrentTime + ".png"
		} else {
			fileName = "MaxDistortionVoters_" + CurrentData + "_" + CurrentTime + ".png"
		}
	} else {
		if isCandidates {
			fileName = "AverageDistortionCandidates_" + CurrentData + "_" + CurrentTime + ".png"
		} else {
			fileName = "AverageDistortionVoters_" + CurrentData + "_" + CurrentTime + ".png"
		}
	}
	if err := p.Save(4*vg.Inch, 4*vg.Inch, fileName); err != nil {
		panic(err)
	}
	return fileName

}

// Parses a set of multiresults to get the unique values of candidates or voters
func MultiParse(results MultiResults, isCandidates bool) []int {
	seenValues := make(map[int]bool)
	for key := range results.Results {
		numEach := strings.Split(key, "_")
		numCandidates := numEach[0]
		numVoters := numEach[1]
		if isCandidates {
			value, _ := strconv.Atoi(string(numCandidates))
			if seenValues[value] {
				continue
			}
			seenValues[value] = true
		} else {
			value, _ := strconv.Atoi(string(numVoters))
			if seenValues[value] {
				continue
			}
			seenValues[value] = true
		}
	}
	keys := make([]int, 0, len(seenValues))
	for key := range seenValues {
		keys = append(keys, key)
	}
	return keys
}

// Parses a set of multiresults to get the distortion values for a specific number of candidates or voters
func ParseCandidatesOrVoters(results MultiResults, numCandiatesOrVoters int, isCandidates bool, isMax bool) map[int]float64 {
	parsedData := make(map[int]float64)
	for key, value := range results.Results {
		numEach := strings.Split(key, "_")
		numCandidates := numEach[0]
		numVoters := numEach[1]
		string_value := 0
		err := error(nil)
		if isCandidates {
			string_value, err = strconv.Atoi(string(numCandidates))
			if err != nil {
				fmt.Println("Error converting string to int:", err)
			}
		} else {
			string_value, err = strconv.Atoi(string(numVoters))
			if err != nil {
				fmt.Println("Error converting string to int:", err)
			}
		}

		if string_value == numCandiatesOrVoters {
			if isMax {
				if isCandidates {
					position, _ := strconv.Atoi(numVoters)
					parsedData[position] = value.MaxDistortion
				} else {
					position, _ := strconv.Atoi(numCandidates)
					parsedData[position] = value.MaxDistortion
				}
			} else {
				if isCandidates {
					position, _ := strconv.Atoi(numVoters)
					parsedData[position] = value.AverageDistortion
				} else {
					position, _ := strconv.Atoi(numCandidates)
					parsedData[position] = value.AverageDistortion
				}
			}
		}
	}
	//fmt.Println("Parsed Data:", parsedData)
	return parsedData
}

// Reads a JSON file and unmarshals the data into a MultiResults struct
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

// Function to graph the results of a single simulation
func GraphOneResults(mapData map[int]float64) {
	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	pts := plotter.XYs{}
	//Append pt
	for key, value := range mapData {
		pts = append(pts, plotter.XY{X: float64(key), Y: value})
	}

	// Add points to the plot as a scatter plot
	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	p.Add(s)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}
