package main

/*
* This file contains the main function to run the voting functions GUI
 */
import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"

	primary "AGT_Midterm/src"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Main function to create the main window and call the other needed functions
func main() {
	runtime.LockOSThread() // Ensure we're locking to the OS main thread

	myApp := app.New()

	mainWindow := myApp.NewWindow("Voting Functions")
	mainWindow.Resize(fyne.NewSize(800, 600))

	titleLabel := widget.NewLabelWithStyle("Voting Functions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	button1 := widget.NewButton("Run Specific Instance", func() {
		selectFiles(myApp)
	})

	button2 := widget.NewButton("Run Random Simulation", func() {
		displaySimulatorVotes(myApp)

	})

	button3 := widget.NewButton("Run Multi Simulation", func() {
		multiSimulation(myApp)
	})

	mainWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		titleLabel,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), button1, button2, button3, layout.NewSpacer()),
		layout.NewSpacer(),
	))

	mainWindow.ShowAndRun()
}

// Function to make a window and has to drop downs to allow the selection of a candidate and voter file
func selectFiles(myApp fyne.App) {
	listFiles := func(dirPath string) ([]string, error) {
		var files []string
		fileInfo, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}
		for _, file := range fileInfo {
			files = append(files, file.Name())
		}
		return files, nil
	}

	// Get candidate and voter files
	candidateFiles, err := listFiles("./Jsons/Candidates")
	if err != nil {
		log.Fatalf("Failed to list candidate files: %v", err)
	}
	voterFiles, err := listFiles("./Jsons/Voters")
	if err != nil {
		log.Fatalf("Failed to list voter files: %v", err)
	}

	candidateDropdown := widget.NewSelect(candidateFiles, nil)
	voterDropdown := widget.NewSelect(voterFiles, nil)
	selectedCandidateFilePath := ""
	selectedVoterFilePath := ""

	candidateDropdown.OnChanged = func(selected string) {
		selectedCandidateFilePath = filepath.Join("./Jsons/Candidates", selected)

	}

	voterDropdown.OnChanged = func(selected string) {
		selectedVoterFilePath = filepath.Join("./Jsons/Voters", selected)

	}

	if len(candidateFiles) > 0 {
		candidateDropdown.SetSelected(candidateFiles[0])
	}
	if len(voterFiles) > 0 {
		voterDropdown.SetSelected(voterFiles[0])
	}

	//Make new button to submit
	submitButton := widget.NewButton("Submit", func() {
		displayVotingResults(myApp, selectedCandidateFilePath, selectedVoterFilePath)
	})
	content := container.NewVBox(
		widget.NewLabel("Select Candidates JSON:"),
		candidateDropdown,
		widget.NewLabel("Select Voters JSON:"),
		voterDropdown,
		submitButton,
	)

	window := myApp.NewWindow("Voting Results")
	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(content)
	window.Show()

}

// Function to display the results of the voting systems
// Creates a new window and displays the results in a table along with candidate and voter positions
func displayVotingResults(myApp fyne.App, candidatesFileName, votersFileName string) {
	// Load candidates and voters from files
	candidateData, err := voting_systems.ReadFromFile(candidatesFileName, "Candidate")
	if err != nil {
		fmt.Println("Error reading candidates file")
		return
	}
	candidates := candidateData.([]voting_systems.Candidate)

	voterData, err := voting_systems.ReadFromFile(votersFileName, "Voter")
	if err != nil {
		fmt.Println("Error reading voters file")
		return
	}
	voters := voterData.([]voting_systems.Voter)

	// Run voting systems and get results
	optimalCost, optCanidate := voting_systems.DetermineOptimalCanidate(append([]voting_systems.Candidate(nil), candidates...), voters)
	stvWinner, stvCanidates, stvRounds := voting_systems.InitiateSTV(append([]voting_systems.Candidate(nil), candidates...), voters)
	//Get STVWinner distortion
	stvWinnerDistortion := voting_systems.GetDistortion(voting_systems.GetSocailCost(stvWinner, voters), optimalCost)
	bordaWinner, bordaCanidates := voting_systems.CalculateBordaWinner(append([]voting_systems.Candidate(nil), candidates...), voters)
	//Get BordaWinner distortion
	bordaWinnerDistortion := voting_systems.GetDistortion(voting_systems.GetSocailCost(bordaWinner, voters), optimalCost)
	pluralityWinner, pluralityCanidates := voting_systems.InitiatePlurality(append([]voting_systems.Candidate(nil), candidates...), voters)
	//Get PluralityWinner distortion
	pluralityWinnerDistortion := voting_systems.GetDistortion(voting_systems.GetSocailCost(pluralityWinner, voters), optimalCost)
	copelandWinner, copelandCanidates := voting_systems.DetermineCopelandWinner(append([]voting_systems.Candidate(nil), candidates...), voters)
	//Get CopelandWinner distortion
	copelandWinnerDistortion := voting_systems.GetDistortion(voting_systems.GetSocailCost(copelandWinner, voters), optimalCost)
	pluralityVetoWinner, vetoCanidates, vetoRounds := voting_systems.InitiatePluralityVeto(append([]voting_systems.Candidate(nil), candidates...), voters)
	//Get PluralityVetoWinner distortion
	pluralityVetoWinnerDistortion := voting_systems.GetDistortion(voting_systems.GetSocailCost(pluralityVetoWinner, voters), optimalCost)

	// Create widgets to display results
	optimalCostLabel := widget.NewLabel(fmt.Sprintf("Optimal Canidate w/ cost: %s %.4f", optCanidate.Name, optimalCost))
	stvWinnerLabel := widget.NewLabel(fmt.Sprintf("STV Winner: %s, distortion: %.4f", stvWinner.Name, stvWinnerDistortion))
	bordaWinnerLabel := widget.NewLabel(fmt.Sprintf("Borda Winner: %s, distortion: %.4f", bordaWinner.Name, bordaWinnerDistortion))
	pluralityWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Winner: %s, distortion: %.4f", pluralityWinner.Name, pluralityWinnerDistortion))
	copelandWinnerLabel := widget.NewLabel(fmt.Sprintf("Copeland Winner: %s, distortion: %.4f", copelandWinner.Name, copelandWinnerDistortion))
	pluralityVetoWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Veto Winner: %s, distortion: %.4f", pluralityVetoWinner.Name, pluralityVetoWinnerDistortion))

	voterTable := container.NewVScroll(createVoterTable(voters))
	voterTable.SetMinSize(fyne.NewSize(400, 200))

	options := []string{"Default Candidates", "STV Candidates", "Borda Candidates", "Plurality Candidates", "Copeland Candidates", "Veto Candidates"}
	candidateArrays := map[string][]voting_systems.Candidate{
		"Default Candidates":   candidates,
		"STV Candidates":       stvCanidates,
		"Borda Candidates":     bordaCanidates,
		"Plurality Candidates": pluralityCanidates,
		"Copeland Candidates":  copelandCanidates,
		"Veto Candidates":      vetoCanidates,
	}

	candidateTable := container.NewVScroll(createCandidateTable(candidates))
	candidateTable.SetMinSize(fyne.NewSize(400, 200))
	roundsDropdown := widget.NewSelect([]string{}, nil)
	roundsDropdown.Hide()

	updateCandidateTable := func(value string) {
		var currentSelectionType string

		if value == "STV Candidates" || value == "Veto Candidates" {
			currentSelectionType = value

			var rounds [][]voting_systems.Candidate
			if value == "STV Candidates" {
				rounds = stvRounds
			} else { // "Veto Candidates"
				rounds = vetoRounds
			}

			roundOptions := make([]string, len(rounds))
			for i := range rounds {
				roundOptions[i] = fmt.Sprintf("Round %d", i+1)
			}
			roundsDropdown.Options = roundOptions
			roundsDropdown.Selected = "" // Reset selection
			roundsDropdown.Refresh()
			roundsDropdown.Show()

			roundsDropdown.OnChanged = func(selected string) {
				var selectedRound int
				fmt.Sscanf(selected, "Round %d", &selectedRound)
				selectedRound -= 1
				var roundCandidates []voting_systems.Candidate
				if currentSelectionType == "STV Candidates" {
					roundCandidates = stvRounds[selectedRound]
				} else { // "Veto Candidates"
					roundCandidates = vetoRounds[selectedRound]
				}

				// Update the table
				newTable := container.NewVScroll(createCandidateTable(roundCandidates))
				candidateTable.Content = newTable.Content
				candidateTable.Refresh()
			}

			if len(roundOptions) > 0 {
				roundsDropdown.SetSelected(roundOptions[0])
			}
		} else { //No rounds
			roundsDropdown.Hide()
			candidateArray := candidateArrays[value]
			newTable := container.NewVScroll(createCandidateTable(candidateArray))
			candidateTable.Content = newTable.Content
			candidateTable.Refresh()
		}
	}

	// Function to update the rounds dropdown based on the selected candidate group
	updateRoundsDropdown := func(candidateGroup string, rounds [][]voting_systems.Candidate) {
		if candidateGroup == "STV Candidates" || candidateGroup == "Veto Candidates" {
			roundOptions := make([]string, len(rounds))
			for i := range rounds {
				roundOptions[i] = fmt.Sprintf("Round %d", i+1)
			}
			roundsDropdown.Options = roundOptions
			roundsDropdown.Selected = roundOptions[0]
			roundsDropdown.Refresh()
			roundsDropdown.Show()
		} else {
			roundsDropdown.Hide()
		}
	}

	dropdown := widget.NewSelect(options, func(selected string) {
		if selected == "STV Candidates" {
			updateRoundsDropdown(selected, stvRounds)
		} else if selected == "Veto Candidates" {
			updateRoundsDropdown(selected, vetoRounds)
		} else {
			roundsDropdown.Hide()
		}
		updateCandidateTable(selected)
	})
	dropdown.PlaceHolder = "Select Candidate Group"

	content := container.NewVBox(
		dropdown,
		roundsDropdown,
		optimalCostLabel,
		stvWinnerLabel,
		bordaWinnerLabel,
		pluralityWinnerLabel,
		copelandWinnerLabel,
		pluralityVetoWinnerLabel,
		widget.NewLabel("Voter Positions:"),
		voterTable,
		widget.NewLabel("Candidate Positions and Votes:"),
		candidateTable,
	)

	window := myApp.NewWindow("Voting Results")
	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(content)

	window.Show()
}

// createVoterTable creates a table widget to display voter positions
func createVoterTable(voters []voting_systems.Voter) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(voters), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			voter := voters[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(voter.Name)
			case 1:
				o.(*widget.Label).SetText(fmt.Sprintf("%.2f", voter.Position))
			case 2:
				o.(*widget.Label).SetText(fmt.Sprintf("%d", voter.Number))
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 200)
	return table
}

// createCandidateTable creates a table widget to display candidate positions, votes, and names
func createCandidateTable(candidates []voting_systems.Candidate) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(candidates), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			candidate := candidates[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(candidate.Name)
			case 1:
				o.(*widget.Label).SetText(fmt.Sprintf("%.2f", candidate.Position))
			case 2:
				o.(*widget.Label).SetText(fmt.Sprintf("%d", candidate.NumVotes))
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 200)
	return table
}

// Function to handle the window for the simulator
// Creates a window with entry widgets for parameters and buttons to run the simulation
func displaySimulatorVotes(app fyne.App) {
	myApp := app

	mainWindow := myApp.NewWindow("Voting Functions")
	mainWindow.Resize(fyne.NewSize(800, 600))

	titleLabel := widget.NewLabelWithStyle("Voting Functions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create entry widgets for parameters with default values
	numRunsEntry := widget.NewEntry()
	numRunsEntry.SetText("10") // Default value
	numRunsEntry.SetPlaceHolder("Number of Runs")

	numCandidatesEntry := widget.NewEntry()
	numCandidatesEntry.SetText("5") // Default value
	numCandidatesEntry.SetPlaceHolder("Number of Candidates")

	numVoterssEntry := widget.NewEntry()
	numVoterssEntry.SetText("5") // Default value
	numVoterssEntry.SetPlaceHolder("Number of Voters")

	maxPositionEntry := widget.NewEntry()
	maxPositionEntry.SetText("1.0") // Default value
	maxPositionEntry.SetPlaceHolder("Max Position")

	minPositionEntry := widget.NewEntry()
	minPositionEntry.SetText("0.0") // Default value
	minPositionEntry.SetPlaceHolder("Min Position")

	totalVotersEntry := widget.NewEntry()
	totalVotersEntry.SetText("100") // Default value
	totalVotersEntry.SetPlaceHolder("Total Voters")

	resultsLabel := widget.NewLabel("")

	// Function to parse entry inputs and run simulation
	candidateFileName, voterFileName := "candidates.json", "voters.json"
	runSimulation := func(system string) {
		numRuns, _ := strconv.Atoi(numRunsEntry.Text)
		numCandidates, _ := strconv.Atoi(numCandidatesEntry.Text)
		nunVoters, _ := strconv.Atoi(numVoterssEntry.Text)
		maxPosition, _ := strconv.ParseFloat(maxPositionEntry.Text, 64)
		minPosition, _ := strconv.ParseFloat(minPositionEntry.Text, 64)
		totalVoters, _ := strconv.Atoi(totalVotersEntry.Text)
		result := ""

		result, candidateFileName, voterFileName = runAndDisplayResults(numRuns, numCandidates, nunVoters, maxPosition, minPosition, totalVoters, system)
		resultsLabel.SetText(result)
	}

	// Create buttons for each voting system
	votingSystems := []string{"STV", "Borda Count", "Plurality", "Copeland", "Plurality Veto"}
	var votingSystemButtons []*widget.Button
	for _, system := range votingSystems {
		button := widget.NewButton(system, func(sys string) func() {
			return func() {
				runSimulation(sys)
			}
		}(system))
		votingSystemButtons = append(votingSystemButtons, button)
	}

	var canvasButtons []fyne.CanvasObject
	for _, button := range votingSystemButtons {
		canvasButtons = append(canvasButtons, fyne.CanvasObject(button))
	}
	var viewDetailsButton = widget.NewButton("View Details", func() {
		candidateFilePath := filepath.Join(".", "Jsons", "Candidates", candidateFileName)
		voterFilePath := filepath.Join(".", "Jsons", "Voters", voterFileName)

		displayVotingResults(myApp, candidateFilePath, voterFilePath)
	})

	// Layout for parameter entries
	paramsContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Number of Runs:"), numRunsEntry,
		widget.NewLabel("Number of Candidates:"), numCandidatesEntry,
		widget.NewLabel("Number of Voters:"), numVoterssEntry,
		widget.NewLabel("Max Position:"), maxPositionEntry,
		widget.NewLabel("Min Position:"), minPositionEntry,
		widget.NewLabel("Total Voters:"), totalVotersEntry,
	)

	mainWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		titleLabel,
		viewDetailsButton,
		paramsContainer,
		container.New(layout.NewGridLayoutWithRows(3), canvasButtons...),
		resultsLabel,
		layout.NewSpacer(),
	))
	mainWindow.Show()
}

// Function to run a simulation and return the results
func runAndDisplayResults(numRuns int, numCandidates int, numVoters int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) (string, string, string) {
	result := ""
	candidateFileName, voterFileName := "candidates.json", "voters.json"
	var maxDistortion float64
	var averageDistortion float64

	switch votingSystem {
	case "STV":
		result, candidateFileName, voterFileName, maxDistortion, averageDistortion = primary.RunScenario(numRuns, numCandidates, numVoters, maxPosition, minPosition, totalVoters, "STV")
	case "Borda Count":
		result, candidateFileName, voterFileName, maxDistortion, averageDistortion = primary.RunScenario(numRuns, numCandidates, numVoters, maxPosition, minPosition, totalVoters, "Borda Count")
	case "Plurality":
		result, candidateFileName, voterFileName, maxDistortion, averageDistortion = primary.RunScenario(numRuns, numCandidates, numVoters, maxPosition, minPosition, totalVoters, "Plurality")
	case "Copeland":
		result, candidateFileName, voterFileName, maxDistortion, averageDistortion = primary.RunScenario(numRuns, numCandidates, numVoters, maxPosition, minPosition, totalVoters, "Copeland")
	case "Plurality Veto":
		result, candidateFileName, voterFileName, maxDistortion, averageDistortion = primary.RunScenario(numRuns, numCandidates, numVoters, maxPosition, minPosition, totalVoters, "Plurality Veto")
	}
	fmt.Println("Max Distortion: ", maxDistortion)
	fmt.Println("Average Distortion: ", averageDistortion)
	return strings.TrimSpace(votingSystem + " Result:\n" + result + "\nAverage Distortion: " + strconv.FormatFloat(averageDistortion, 'f', -1, 64) + "\n\n"), candidateFileName, voterFileName

}

// Handles the multi-simulation window
func multiSimulation(app fyne.App) {
	myApp := app

	mainWindow := myApp.NewWindow("Voting Functions")
	mainWindow.Resize(fyne.NewSize(800, 800))

	titleLabel := widget.NewLabelWithStyle("Voting Functions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create entry widgets for parameters with default values
	numRunsEntry := widget.NewEntry()
	numRunsEntry.SetText("10") // Default value
	numRunsEntry.SetPlaceHolder("Number of Runs")

	minNumCandidatesEntry := widget.NewEntry()
	minNumCandidatesEntry.SetText("1") // Default value
	minNumCandidatesEntry.SetPlaceHolder("Min Number of Candidates")

	maxNumCandidatesEntry := widget.NewEntry()
	maxNumCandidatesEntry.SetText("5") // Default value
	maxNumCandidatesEntry.SetPlaceHolder("Max Number of Candidates")

	minNumVoterssEntry := widget.NewEntry()
	minNumVoterssEntry.SetText("1") // Default value
	minNumVoterssEntry.SetPlaceHolder("Min Number of Voters")

	maxNumVoterssEntry := widget.NewEntry()
	maxNumVoterssEntry.SetText("5") // Default value
	maxNumVoterssEntry.SetPlaceHolder("Max Number of Voters")

	maxPositionEntry := widget.NewEntry()
	maxPositionEntry.SetText("1.0") // Default value
	maxPositionEntry.SetPlaceHolder("Max Position")

	minPositionEntry := widget.NewEntry()
	minPositionEntry.SetText("0.0") // Default value
	minPositionEntry.SetPlaceHolder("Min Position")

	totalVotersEntry := widget.NewEntry()
	totalVotersEntry.SetText("100") // Default value
	totalVotersEntry.SetPlaceHolder("Total Voters")

	resultsLabel := widget.NewMultiLineEntry()
	resultsLabel.Wrapping = fyne.TextWrapWord
	//resultsLabel.Disable()
	scrollContainer := container.NewScroll(resultsLabel)
	scrollContainer.SetMinSize(fyne.NewSize(200, 200))

	fileName := ""
	minNumCandidates := 0
	maxNumCandidates := 0
	minNumVoters := 0
	maxNumVoters := 0
	candidateOptions := []string{"All Candidates"}
	voterOptions := []string{"All Voters"}

	candidateDropdown := widget.NewSelect(candidateOptions, nil)
	voterDropdown := widget.NewSelect(voterOptions, nil)
	selectedNumCandidates := "All Candidates"
	selectedNumVoters := "All Voters"

	if len(candidateOptions) > 0 {
		candidateDropdown.SetSelected(candidateOptions[0])
	}
	if len(voterOptions) > 0 {
		voterDropdown.SetSelected(voterOptions[0])
	}
	resultsLabel.SetText(selectedNumCandidates + " " + selectedNumVoters)
	result := make(map[primary.MultiInput]primary.MultiOutput)

	updateResults := func() {
		var keys []primary.MultiInput
		for k := range result {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			if keys[i].NumCandidates == keys[j].NumCandidates {
				return keys[i].NumVoters < keys[j].NumVoters
			}
			return keys[i].NumCandidates < keys[j].NumCandidates
		})

		outputText := ""
		for _, key := range keys {
			value := result[key]
			candidateSelectedNumberString := strings.Split(selectedNumCandidates, " ")
			candidateSelectedNumber := candidateSelectedNumberString[0]
			voterSelectedNumberString := strings.Split(selectedNumVoters, " ")
			voterSelectedNumber := voterSelectedNumberString[0]
			if candidateSelectedNumber == "A" {
				candidateSelectedNumber = "All"
			}
			if voterSelectedNumber == "A" {
				voterSelectedNumber = "All"
			}
			candidateMatch := strconv.Itoa(key.NumCandidates) == candidateSelectedNumber || candidateSelectedNumber == "All"
			voterMatch := strconv.Itoa(key.NumVoters) == voterSelectedNumber || voterSelectedNumber == "All"

			if candidateMatch && voterMatch {
				outputText += fmt.Sprintf("Number of Candidates: %d, Number of Voters: %d\n", key.NumCandidates, key.NumVoters) + value.Result + "\nMax Distortion: " + fmt.Sprintf("%.2f", value.MaxDistortion) + ", Average Distortion: " + fmt.Sprintf("%.2f", value.AverageDistortion) + "\n\n"
			}
		}
		resultsLabel.SetText(outputText)
	}

	candidateDropdown.SetSelected(selectedNumCandidates)
	voterDropdown.SetSelected(selectedNumVoters)

	candidateDropdown.OnChanged = func(selected string) {
		selectedNumCandidates = selected
		updateResults()
	}

	voterDropdown.OnChanged = func(selected string) {
		selectedNumVoters = selected
		updateResults()
	}

	runSimulation := func(system string) {
		numRuns, _ := strconv.Atoi(numRunsEntry.Text)
		minNumCandidates, _ = strconv.Atoi(minNumCandidatesEntry.Text)
		maxNumCandidates, _ = strconv.Atoi(maxNumCandidatesEntry.Text)
		minNumVoters, _ = strconv.Atoi(minNumVoterssEntry.Text)
		maxNumVoters, _ = strconv.Atoi(maxNumVoterssEntry.Text)
		maxPosition, _ := strconv.ParseFloat(maxPositionEntry.Text, 64)
		minPosition, _ := strconv.ParseFloat(minPositionEntry.Text, 64)
		totalVoters, _ := strconv.Atoi(totalVotersEntry.Text)

		for i := minNumCandidates; i <= maxNumCandidates; i++ {
			string_option := strconv.Itoa(i) + " Candidates"
			candidateOptions = append(candidateOptions, string_option)
		}
		for i := minNumVoters; i <= maxNumVoters; i++ {
			string_option := strconv.Itoa(i) + " Voters"
			voterOptions = append(voterOptions, string_option)
		}

		candidateDropdown.SetOptions(candidateOptions)
		voterDropdown.SetOptions(voterOptions)

		result, _, _, fileName = multiSimResults(app, minNumCandidates, maxNumCandidates, minNumVoters, maxNumVoters, numRuns, maxPosition, minPosition, totalVoters, system)
		updateResults()
	}

	// Create buttons for each voting system
	votingSystems := []string{"STV", "Borda Count", "Plurality", "Copeland", "Plurality Veto"}
	var votingSystemButtons []*widget.Button
	for _, system := range votingSystems {
		button := widget.NewButton(system, func(sys string) func() {
			return func() {
				runSimulation(sys)
			}
		}(system))
		votingSystemButtons = append(votingSystemButtons, button)
	}

	var graphsButton = widget.NewButton("Create Graphs", func() {
		DisplayGraphs(app, fileName, true, true)
	})

	fmt.Println(selectedNumCandidates)
	fmt.Println(selectedNumVoters)

	var canvasButtons []fyne.CanvasObject
	for _, button := range votingSystemButtons {
		canvasButtons = append(canvasButtons, fyne.CanvasObject(button))
	}

	// Layout for parameter entries
	paramsContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Number of Runs:"), numRunsEntry,
		widget.NewLabel("Min Number of Candidates:"), minNumCandidatesEntry,
		widget.NewLabel("Max Number of Candidates:"), maxNumCandidatesEntry,
		widget.NewLabel("Min Number of Voters:"), minNumVoterssEntry,
		widget.NewLabel("Max Number of Voters:"), maxNumVoterssEntry,
		widget.NewLabel("Max Position:"), maxPositionEntry,
		widget.NewLabel("Min Position:"), minPositionEntry,
		widget.NewLabel("Total Voters:"), totalVotersEntry,
	)

	mainWindow.SetContent(container.NewVBox(
		graphsButton,
		layout.NewSpacer(),
		titleLabel,
		paramsContainer,
		container.New(layout.NewGridLayoutWithRows(3), canvasButtons...),
		container.New(layout.NewGridLayoutWithColumns(2), candidateDropdown, voterDropdown),
		scrollContainer,
		layout.NewSpacer(),
	))
	mainWindow.Show()

}

func multiSimResults(app fyne.App, minCandidates int, maxCandidates int, minVoters int, maxVoters int, numRuns int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) (map[primary.MultiInput]primary.MultiOutput, string, string, string) {
	result := ""
	multiResults := make(map[primary.MultiInput]primary.MultiOutput)
	fileName := ""

	switch votingSystem {
	case "STV":
		multiResults, fileName = primary.Multi_sim(minCandidates, maxCandidates, minVoters, maxVoters, numRuns, maxPosition, minPosition, totalVoters, "STV")
	case "Borda Count":
		multiResults, fileName = primary.Multi_sim(minCandidates, maxCandidates, minVoters, maxVoters, numRuns, maxPosition, minPosition, totalVoters, "Borda Count")
	case "Plurality":
		multiResults, fileName = primary.Multi_sim(minCandidates, maxCandidates, minVoters, maxVoters, numRuns, maxPosition, minPosition, totalVoters, "Plurality")
	case "Copeland":
		multiResults, fileName = primary.Multi_sim(minCandidates, maxCandidates, minVoters, maxVoters, numRuns, maxPosition, minPosition, totalVoters, "Copeland")
	case "Plurality Veto":
		multiResults, fileName = primary.Multi_sim(minCandidates, maxCandidates, minVoters, maxVoters, numRuns, maxPosition, minPosition, totalVoters, "Plurality Veto")
	}
	for key, value := range multiResults {
		result += fmt.Sprintf("Number of Candidates: %d, Number of Voters: %d\n", key.NumCandidates, key.NumVoters)
		result += value.Result + "\n"
		result += fmt.Sprintf("Max Distortion: %.2f, Average Distortion: %.2f\n\n", value.MaxDistortion, value.AverageDistortion)
	}
	//fmt.Println(result)
	return multiResults, "candidates.json", "voters.json", fileName

}

// Function to display/save the graphs
func DisplayGraphs(app fyne.App, fileName string, isMax bool, isCandidates bool) {
	myApp := app
	confirmLabelText := "Select a graph to generate"
	confirmLabel := widget.NewLabel(confirmLabelText)

	mainWindow := myApp.NewWindow("Graph Generator")
	mainWindow.Resize(fyne.NewSize(800, 600))

	titleLabel := widget.NewLabelWithStyle("Graph Generator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	CandidateMaxGraphButton := widget.NewButton("Graph Candidates Max Distortion", func() {
		primary.ReadAndGraphMultiResults(fileName, true, true)
		confirmLabelText = "Max Candidate Distortion Graph Created"
		confirmLabel.SetText(confirmLabelText)
	})

	CandidateAvgGraphButton := widget.NewButton("Graph Candidates Average Distortion", func() {
		primary.ReadAndGraphMultiResults(fileName, false, true)
		confirmLabelText = "Average Candidate Distortion Graph Created"
		confirmLabel.SetText(confirmLabelText)
	})

	VoterMaxGraphButton := widget.NewButton("Graph Voters Max Distortion", func() {
		primary.ReadAndGraphMultiResults(fileName, true, false)
		confirmLabelText = "Max Voter Distortion Graph Created"
		confirmLabel.SetText(confirmLabelText)
	})

	VoterAvgGraphButton := widget.NewButton("Graph Voters Average Distortion", func() {
		primary.ReadAndGraphMultiResults(fileName, false, false)
		confirmLabelText = "Average Voter Distortion Graph Created"
		confirmLabel.SetText(confirmLabelText)
	})

	//Create label

	mainWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		titleLabel,
		container.New(layout.NewGridLayoutWithRows(2), CandidateMaxGraphButton, CandidateAvgGraphButton, VoterMaxGraphButton, VoterAvgGraphButton),
		layout.NewSpacer(),
		confirmLabel,
	))
	mainWindow.Show()
}
