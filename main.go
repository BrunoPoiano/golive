package main

import (
	"fmt"
	interfaces "main/interface"
	"main/models"
	pw "main/pw_functions"
	"math"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

var program *tea.Program
var volumeRate = 0.01

type MainModel struct {
	models.MainModel
}

func initialModel() MainModel {
	lists := pw.RefresLists(models.MainModel{})

	return MainModel{
		MainModel: models.MainModel{
			Padding: 2,
			Cursor:  0,
			Input: models.Input{
				Items:  lists.Input.Items,
				Volume: 1.0,
			},
			Output: models.Output{
				Items:  lists.Output.Items,
				Volume: 1.0,
			},
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case models.LevelMsg:
		m.Level.Value = string(msg)

	// key press actions
	case tea.KeyPressMsg:
		switch msg.String() {

		//Actions
		case "r":
			m.MainModel = pw.RefresLists(m.MainModel)
			return m, nil

		// Play the currently setup
		case "p":
			if m.Play == nil {
				m.Play = pw.Play(m.MainModel)
				m.Debug = fmt.Sprintf("%d", m.Play.Process.Pid)
				m.Level.Process = pw.MonitorChanel(program, m.Input.Items[m.Input.Selected].Info.Props.NodeName)
			}
			return m, nil

		// exit the program.
		case "ctrl+c", "q":
			pw.KillProcesses(m.MainModel)
			return m, tea.Quit

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.Play == nil && m.Cursor < (len(m.Input.Items)-1+len(m.Output.Items)) {
				m.Cursor++
			}
			return m, nil

		case "left":
			if m.Output.Volume > 0 {
				m.Output.Volume = math.Max(0, m.Output.Volume-volumeRate)
			}
			if m.Play != nil {
				id := m.Output.Items[m.Output.Selected].Info.Props.NodeId
				go pw.ChangeVolume(id, m.Output.Volume)
			}
			return m, nil
		case "right":
			if m.Output.Volume < 1.0 {
				m.Output.Volume += volumeRate
			}
			if m.Play != nil {
				id := m.Output.Items[m.Output.Selected].Info.Props.NodeId
				go pw.ChangeVolume(id, m.Output.Volume)
			}
			return m, nil

		case "a":
			if m.Input.Volume > 0 {
				m.Input.Volume = math.Max(0, m.Input.Volume-volumeRate)
			}
			if m.Play != nil {
				id := m.Input.Items[m.Input.Selected].Id
				go pw.ChangeVolume(id, m.Input.Volume)
			}
			return m, nil
		case "d":
			if m.Input.Volume < 1.0 {
				m.Input.Volume += volumeRate
			}
			if m.Play != nil {
				id := m.Input.Items[m.Input.Selected].Id
				go pw.ChangeVolume(id, m.Input.Volume)
			}
			return m, nil

			//Interactions
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.Play == nil && m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		// kill the currently proccess of playing
		case "x":
			m.MainModel = pw.KillProcesses(m.MainModel)
			return m, nil

		// The "enter" key and the space bar toggle the selected state
		case "enter", "space":
			if m.Play == nil {
				if m.Cursor < len(m.Input.Items) {
					m.Input.Selected = m.Cursor
				} else {
					m.Output.Selected = m.Cursor - len(m.Input.Items)
				}
			}
			return m, nil

		}
	}

	// Return the updated Inputs to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainModel) View() tea.View {

	var view strings.Builder

	view.WriteString(interfaces.Header())
	view.WriteString("\n")
	if m.Play != nil {
		view.WriteString(interfaces.Playing(m.MainModel))
	} else {
		view.WriteString(interfaces.ListItems(m.MainModel))
	}
	view.WriteString("\n   a:  decrease input volume |     d: increase input volume")
	view.WriteString("\nleft: decrease output volume | right: increase output volume")

	if m.Play != nil {
		view.WriteString("\nx: Stop | q: quit")
	} else {
		view.WriteString("\np: play | r: refresh lists | q: quit")
	}

	viewBorder := interfaces.Border(m.Padding, m.Width).Render(view.String())

	screen := tea.NewView(viewBorder)
	screen.AltScreen = true
	return screen
}

func main() {
	program = tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
