package main

import (
	"fmt"
	interfaces "main/interface"
	"main/models"
	pw "main/pw_functions"
	"os"

	tea "charm.land/bubbletea/v2"
)

var program *tea.Program

type MainModel struct {
	models.MainModel
}

func initialModel() MainModel {

	inputsList, err := pw.ReturnList(models.InputList)
	if err != nil {
		panic(err.Error())
	}
	outputList, err := pw.ReturnList(models.OutputList)
	if err != nil {
		panic(err.Error())
	}

	return MainModel{
		MainModel: models.MainModel{
			Cursor: 0,
			Input: models.Input{
				Items: inputsList,
			},
			Output: models.Output{
				Items: outputList,
			},
		},
	}
}
func refresLists(m MainModel) MainModel {
	inputsList, err := pw.ReturnList(models.InputList)
	if err != nil {
		panic(err.Error())
	}
	outputList, err := pw.ReturnList(models.OutputList)
	if err != nil {
		panic(err.Error())
	}

	m.Input.Items = inputsList
	m.Output.Items = outputList

	return m
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case models.LevelMsg:
		m.Level = string(msg)

	// key press actions
	case tea.KeyPressMsg:
		switch msg.String() {

		// exit the program.
		case "ctrl+c", "q":
			if m.PlayProcess != nil && m.PlayProcess.Process != nil {
				m.PlayProcess.Process.Kill()
			}
			if m.LevelProcess != nil && m.LevelProcess.Process != nil {
				m.LevelProcess.Process.Kill()
			}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		// kill the currently proccess of playing
		case "s":
			if m.PlayProcess != nil && m.PlayProcess.Process != nil {
				m.PlayProcess.Process.Kill()
				m.PlayProcess = nil
			}
			if m.LevelProcess != nil && m.LevelProcess.Process != nil {
				m.LevelProcess.Process.Kill()
				m.LevelProcess = nil
			}

		case "r":
			m = refresLists(m)
			return m, nil

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.Cursor < (len(m.Input.Items) - 1 + len(m.Output.Items)) {
				m.Cursor++
			}

		// Play the currently setup
		case "p":
			m.PlayProcess = pw.Play(m.Input.Items[m.Input.Selected], m.Output.Items[m.Output.Selected])
			go pw.MonitorChanel(m.LevelProcess, program, m.Input.Items[m.Input.Selected])

		// The "enter" key and the space bar toggle the selected state
		case "enter", "space":
			if m.Cursor < len(m.Input.Items) {
				m.Input.Selected = m.Cursor
			} else {
				m.Output.Selected = m.Cursor - len(m.Input.Items)
			}
		}
	}

	// Return the updated Inputs to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() tea.View {
	if m.PlayProcess != nil {
		return tea.NewView(interfaces.Playing(m.MainModel))
	}
	return tea.NewView(interfaces.ListItems(m.MainModel))
}

func main() {
	program = tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
