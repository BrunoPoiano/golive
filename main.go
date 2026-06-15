package main

import (
	"fmt"
	interfaces "main/interface"
	"main/models"
	pw "main/pw_functions"
	"math"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var program *tea.Program
var maxVolume = 2.0
var volumeRate = 0.01

type MainModel struct {
	models.MainModel
}

func initialModel() MainModel {
	lists, error := pw.RefreshLists(models.MainModel{})

	return MainModel{
		MainModel: models.MainModel{
			Padding: 2,
			Cursor:  0,
			Error:   error,
			Input: models.Input{
				Items: lists.Input.Items,
				Volume: models.Volume{
					Value: 1.0,
					Mute:  false,
				},
			},
			Output: models.Output{
				Items: lists.Output.Items,
				Volume: models.Volume{
					Value: 1.0,
					Mute:  false,
				},
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
		m.Level.PeakLevel = msg.PeakLevel
		m.Level.RMSLevel = msg.RMSLevel

	case models.ErrorMsg:
		m.Error = msg

	// key press actions
	case tea.KeyPressMsg:
		switch msg.String() {

		//Actions
		case "r":
			m.Error = nil
			var err error
			m.MainModel, err = pw.RefreshLists(m.MainModel)
			if err != nil {
				m.Error = err
			}
			return m, nil

		// Play the currently setup
		case "p":
			if m.Play.Cmd == nil && len(m.Input.Items) != 0 && len(m.Output.Items) != 0 {
				m.MainModel = pw.Play(program, m.MainModel)
				m.MainModel = pw.MonitorChannel(program, m.MainModel)
			}
			return m, nil

		// kill the currently proccess of playing
		case "x":
			m.MainModel = pw.KillProcesses(program, m.MainModel)
			return m, nil

		// exit the program.
		case "ctrl+c", "q":
			pw.KillProcesses(program, m.MainModel)
			return m, tea.Quit

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.Play.Cmd == nil && m.Cursor < (len(m.Input.Items)-1+len(m.Output.Items)) {
				m.Cursor++
			}
			return m, nil

		case "m":
			m.Output.Volume.Mute = !m.Output.Volume.Mute
			if m.Play.Cmd != nil {
				id := m.Output.Items[m.Output.Selected].Id
				go pw.ChangeVolume(id, m.Output.Volume)
			}
		case "n":
			m.Input.Volume.Mute = !m.Input.Volume.Mute
			if m.Play.Cmd != nil {
				id := m.Input.Items[m.Input.Selected].Id
				go pw.ChangeVolume(id, m.Input.Volume)
			}
		// output volume
		case "left":
			if m.Output.Volume.Value > 0 {
				m.Output.Volume.Value = math.Max(0, m.Output.Volume.Value-volumeRate)
			}
			if m.Play.Cmd != nil {
				id := m.Output.Items[m.Output.Selected].Id
				go pw.ChangeVolume(id, m.Output.Volume)
			}
			return m, nil

		// output volume
		case "right":
			if m.Output.Volume.Value < maxVolume {
				m.Output.Volume.Value = math.Min(maxVolume, m.Output.Volume.Value+volumeRate)
			}
			if m.Play.Cmd != nil {
				id := m.Output.Items[m.Output.Selected].Id
				go pw.ChangeVolume(id, m.Output.Volume)
			}
			return m, nil

		// Input Volume
		case "a":
			if m.Input.Volume.Value > 0 {
				m.Input.Volume.Value = math.Max(0, m.Input.Volume.Value-volumeRate)
			}
			if m.Play.Cmd != nil {
				id := m.Input.Items[m.Input.Selected].Id
				go pw.ChangeVolume(id, m.Input.Volume)
			}
			return m, nil

		// Input Volume
		case "d":
			if m.Input.Volume.Value < maxVolume {
				m.Input.Volume.Value = math.Min(maxVolume, m.Input.Volume.Value+volumeRate)
			}
			if m.Play.Cmd != nil {
				id := m.Input.Items[m.Input.Selected].Id
				go pw.ChangeVolume(id, m.Input.Volume)
			}
			return m, nil

		//Interactions
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.Play.Cmd == nil && m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		// The "enter" key and the space bar toggle the selected state
		case "enter", "space":
			if m.Play.Cmd == nil {
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

	header, left, right := interfaces.View(m.MainModel)

	left = lipgloss.NewStyle().Width(70).Render(left)

	content := lipgloss.JoinHorizontal(lipgloss.Left, left, right)
	if m.Width <= 110 {
		right = lipgloss.NewStyle().MarginTop(1).Render(right)
		content = lipgloss.JoinVertical(lipgloss.Top, left, right)
	} else {
		right = lipgloss.NewStyle().MarginLeft(2).Render(right)
	}

	finalScreen := lipgloss.JoinVertical(lipgloss.Top, header, content)
	viewBorder := interfaces.Border(m.Padding, m.Width).Render(finalScreen)

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
