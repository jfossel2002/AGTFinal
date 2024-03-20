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
	optimalCost, _ := voting_systems.DetermineOptimalCanidate(candidates, voters)
	stvWinner := voting_systems.InitiateSTV(candidates, len(voters), voters)
	bordaWinner := voting_systems.CalculateBordaWinner(candidates, voters)
	pluralityWinner := voting_systems.InitiatePlurality(candidates, voters)
	copelandWinner := voting_systems.DetermineCopelandWinner(candidates, voters)
	pluralityVetoWinner := voting_systems.InitiatePluralityVeto(candidates, voters)

	// Create widgets to display results
	optimalCostLabel := widget.NewLabel(fmt.Sprintf("Optimal Cost: %.2f", optimalCost))
	stvWinnerLabel := widget.NewLabel(fmt.Sprintf("STV Winner: %s", stvWinner.Name))
	bordaWinnerLabel := widget.NewLabel(fmt.Sprintf("Borda Winner: %s", bordaWinner.Name))
	pluralityWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Winner: %s", pluralityWinner.Name))
	copelandWinnerLabel := widget.NewLabel(fmt.Sprintf("Copeland Winner: %s", copelandWinner.Name))
	pluralityVetoWinnerLabel := widget.NewLabel(fmt.Sprintf("Plurality Veto Winner: %s", pluralityVetoWinner.Name))

	// Create tables to display voter and candidate positions
	voterTable := createVoterTable(voters)
	candidateTable := createCandidateTable(candidates)

	// Create a container to hold all the widgets
	content := container.NewVBox(
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

	// Create a window to hold the content
	window := myApp.NewWindow("Voting Results")
	window.Resize(fyne.NewSize(800, 600))
	window.SetContent(container.New(layout.NewStackLayout(), content))

	// Show the window
	window.Show()
}

// createVoterTable creates a table widget to display voter positions
func createVoterTable(voters []voting_systems.Voter) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(voters), 2
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
