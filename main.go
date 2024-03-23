// File to handle main functions
package main

import (
	voting_systems "AGT_Midterm/src/systems"
	"fmt"

	primary "AGT_Midterm/src"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()

	mainWindow := myApp.NewWindow("Voting Functions")
	mainWindow.Resize(fyne.NewSize(800, 600))

	titleLabel := widget.NewLabelWithStyle("Voting Functions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	// titleLabel.SetTextSize = 20

	button1 := widget.NewButton("Run Specific Instance", func() {
		displayVotingResults(myApp)
	})

	button2 := widget.NewButton("Run Random Simulation", func() {
		displaySimulatorVotes()

		// Implement this function to display the other page
	})

	mainWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		titleLabel,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), button1, button2, layout.NewSpacer()),
		layout.NewSpacer(),
	))

	mainWindow.ShowAndRun()
}

func displayVotingResults(myApp fyne.App) {
	// Load candidates and voters from files
	candidateData, err := voting_systems.ReadFromFile("canidates.json", "Candidate")
	if err != nil {
		// Handle error
	}
	candidates := candidateData.([]voting_systems.Candidate)

	voterData, err := voting_systems.ReadFromFile("voters.json", "Voter")
	if err != nil {
		// Handle error
	}
	voters := voterData.([]voting_systems.Voter)

	// Run voting systems and get results
	optimalCost, optCanidate := voting_systems.DetermineOptimalCanidate(append([]voting_systems.Candidate(nil), candidates...), voters)
	stvWinner, stvCanidates := voting_systems.InitiateSTV(append([]voting_systems.Candidate(nil), candidates...), voters)
	bordaWinner, bordaCanidates := voting_systems.CalculateBordaWinner(append([]voting_systems.Candidate(nil), candidates...), voters)
	pluralityWinner, pluralityCanidates := voting_systems.InitiatePlurality(append([]voting_systems.Candidate(nil), candidates...), voters)
	copelandWinner, copelandCanidates := voting_systems.DetermineCopelandWinner(append([]voting_systems.Candidate(nil), candidates...), voters)
	pluralityVetoWinner, vetoCanidates := voting_systems.InitiatePluralityVeto(append([]voting_systems.Candidate(nil), candidates...), voters)

	// Create widgets to display results
	optimalCostLabel := widget.NewLabel(fmt.Sprintf("Optimal Canidate w/ cost: %s %.2f", optCanidate.Name, optimalCost))
	stvWinnerLabel := widget.NewLabel(fmt.Sprintf("STV Winner: %s", stvWinner.Name))
	bordaWinnerLabel := widget.NewLabel(fmt.Sprintf("Borda Winner: %s", bordaWinner.Name))
	pluralityWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Winner: %s", pluralityWinner.Name))
	copelandWinnerLabel := widget.NewLabel(fmt.Sprintf("Copeland Winner: %s", copelandWinner.Name))
	pluralityVetoWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Veto Winner: %s", pluralityVetoWinner.Name))

	voterTable := container.NewVScroll(createVoterTable(voters))
	voterTable.SetMinSize(fyne.NewSize(400, 200))

	options := []string{"Default Candidates", "STV Candidates", "Borda Candidates", "Plurality Candidates", "Copeland Candidates", "Veto Candidates"}
	candidateArrays := map[string][]voting_systems.Candidate{
		"Default Candidates":   candidates, // Original array
		"STV Candidates":       stvCanidates,
		"Borda Candidates":     bordaCanidates,
		"Plurality Candidates": pluralityCanidates,
		"Copeland Candidates":  copelandCanidates,
		"Veto Candidates":      vetoCanidates,
	}

	// Initially, display the default candidate table
	candidateTable := container.NewVScroll(createCandidateTable(candidates))
	candidateTable.SetMinSize(fyne.NewSize(400, 200))

	// Dropdown selection changed function
	updateCandidateTable := func(value string) {
		candidateArray := candidateArrays[value]
		newTable := container.NewVScroll(createCandidateTable(candidateArray))
		candidateTable.Content = newTable.Content
		candidateTable.Refresh()
	}

	dropdown := widget.NewSelect(options, updateCandidateTable)
	dropdown.PlaceHolder = "Select Candidate Group"

	content := container.NewVBox(
		dropdown, // Add the dropdown to the layout
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

	// Create a window to hold the content with a scrollable container if necessary
	window := myApp.NewWindow("Voting Results")
	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(content)

	// Show the window
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

// createCandidateTable creates a table widget to display candidate positions and votes
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

func displaySimulatorVotes() {
	myApp := app.New()

	mainWindow := myApp.NewWindow("Voting Functions")
	mainWindow.Resize(fyne.NewSize(800, 600))

	titleLabel := widget.NewLabelWithStyle("Voting Functions", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	resultsLabel := widget.NewLabel("")

	// Create buttons for each voting system
	votingSystems := []string{"STV", "Borda Count", "Plurality", "Copeland", "Plurality Veto"}
	var votingSystemButtons []*widget.Button
	for _, system := range votingSystems {
		button := widget.NewButton(system, func(sys string) func() {
			return func() {
				resultsLabel.SetText(runAndDisplayResults(10, 5, 1.0, 0.0, 100, sys))
			}
		}(system))
		votingSystemButtons = append(votingSystemButtons, button)
	}

	// Convert votingSystemButtons to []fyne.CanvasObject
	var canvasButtons []fyne.CanvasObject
	for _, button := range votingSystemButtons {
		canvasButtons = append(canvasButtons, fyne.CanvasObject(button))
	}

	showResultsButton := widget.NewButton("Show Results", func() {
		results := runAndDisplayResults(10, 5, 1.0, 0.0, 100, "STV")
		results += runAndDisplayResults(10, 5, 1.0, 0.0, 100, "Borda Count")
		results += runAndDisplayResults(10, 5, 1.0, 0.0, 100, "Plurality")
		results += runAndDisplayResults(10, 5, 1.0, 0.0, 100, "Copeland")
		results += runAndDisplayResults(10, 5, 1.0, 0.0, 100, "Plurality Veto")
		resultsLabel.SetText(results)
	})

	mainWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		titleLabel,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), showResultsButton, layout.NewSpacer()),
		layout.NewSpacer(),
		container.New(layout.NewGridLayoutWithRows(3), canvasButtons...), // Use GridLayoutWithRows to display multiple rows
		resultsLabel,
		layout.NewSpacer(),
	))
	mainWindow.Show()
}

func runAndDisplayResults(numRuns int, numCandidates int, maxPosition float64, minPosition float64, totalVoters int, votingSystem string) string {
	result := ""

	switch votingSystem {
	case "STV":
		result = primary.RunScenario(numRuns, numCandidates, maxPosition, minPosition, totalVoters, "STV")
	case "Borda Count":
		result = primary.RunScenario(numRuns, numCandidates, maxPosition, minPosition, totalVoters, "Borda Count")
	case "Plurality":
		result = primary.RunScenario(numRuns, numCandidates, maxPosition, minPosition, totalVoters, "Plurality")
	case "Copeland":
		result = primary.RunScenario(numRuns, numCandidates, maxPosition, minPosition, totalVoters, "Copeland")
	case "Plurality Veto":
		result = primary.RunScenario(numRuns, numCandidates, maxPosition, minPosition, totalVoters, "Plurality Veto")
	}
	return strings.TrimSpace(votingSystem + " Result:\n" + result + "\n\n")

}
