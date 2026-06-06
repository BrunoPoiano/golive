package main

import (
	"fmt"
	"main/parser"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

type MainModel struct {
	inputs  Inputs
	output  Output
	process *exec.Cmd
	cursor  int
}

type Inputs struct {
	inputs   []string
	selected int
}

type Output struct {
	output   []string
	selected int
}

func initialModel() MainModel {

	inputsList, err := parser.ReturnInputList()
	if err != nil {
		panic(err.Error())
	}
	outputList, err := parser.ReturnOuputList()
	if err != nil {
		panic(err.Error())
	}

	return MainModel{
		cursor: 0,

		inputs: Inputs{
			inputs: inputsList,
		},
		output: Output{
			output: outputList,
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			if m.process.Process != nil {
				m.process.Process.Kill()
			}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "m":
			if m.process.Process != nil {
				m.process.Process.Kill()
				m.process = nil
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < (len(m.inputs.inputs) - 1 + len(m.output.output)) {
				m.cursor++
			}

		case "p":
			m.process = parser.Play(m.inputs.inputs[m.inputs.selected], m.output.output[m.output.selected])

		// The "enter" key and the space bar toggle the selected state
		// for the item that the cursor is pointing at.
		case "enter", "space":
			if m.cursor < len(m.inputs.inputs) {
				m.inputs.selected = m.cursor
			} else {
				m.output.selected = m.cursor - len(m.inputs.inputs)
			}
		}
	}

	// Return the updated Inputs to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() tea.View {
	s := ""

	if m.process != nil {
		s = "Playing \n\n "

		s += fmt.Sprintf("Input: %s \n", m.inputs.inputs[m.inputs.selected])
		s += fmt.Sprintf("Output: %s \n", m.output.output[m.output.selected])
		s += "\n Press m to Stop \n"

		return tea.NewView(s)
	}

	// The header
	s += "Select Input\n\n"

	// Iterate over our choices
	for i, choice := range m.inputs.inputs {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.inputs.selected == i {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "Select output\n\n"

	for i, choice := range m.output.output {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i+len(m.inputs.inputs) {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.output.selected == i {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress p to play.\n"
	s += "Press q to quit.\n"

	// Send the UI for rendering
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
