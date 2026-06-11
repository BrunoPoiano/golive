package interfaces

import (
	"fmt"
	"main/models"
	"strings"

	"charm.land/lipgloss/v2"
)

func Header() string {

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50C878")).
		Render(".:GoLive:.\n")
}

func Playing(m models.MainModel) string {

	var s strings.Builder
	fmt.Fprintf(&s, "%s\n", generateMeter(m.Level.Value))
	// fmt.Fprintf(&s, "debug %d\n", m.Input.Items[m.Input.Selected].Id)
	s.WriteString("\nPlaying\n")
	fmt.Fprintf(&s, " Input: %d%% | %s\n", int(m.Input.Volume*100/1), (m.Input.Items[m.Input.Selected].Info.Props.NodeDescription))
	fmt.Fprintf(&s, "Output: %d%% | %s\n", int(m.Output.Volume*100/1), (m.Output.Items[m.Output.Selected].Info.Props.NodeDescription))

	return s.String()
}

func ListItems(m models.MainModel) string {

	var s strings.Builder
	// s = fmt.Sprintf("debug: %s\n", m.Debug)
	fmt.Fprintf(&s, "Inputs: %d%%\n", int(m.Input.Volume*100/1))
	// Iterate over our choices
	for i, choice := range m.Input.Items {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.Input.Selected == i {
			checked = "x" // selected!
		}

		// Render the row
		fmt.Fprintf(&s, "%s [%s] %s\n", cursor, checked, choice.Info.Props.NodeDescription)
	}

	fmt.Fprintf(&s, "\nOutputs: %d%%\n", int(m.Output.Volume*100/1))

	for i, choice := range m.Output.Items {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.Cursor == i+len(m.Input.Items) {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.Output.Selected == i {
			checked = "x" // selected!
		}

		// Render the row
		fmt.Fprintf(&s, "%s [%s] %s\n", cursor, checked, choice.Info.Props.NodeDescription)
	}

	return s.String()
}

func WidthCalc(m models.MainModel, v_width float64) int {
	width := (float64(m.Width) * v_width) - float64(m.Padding)
	return int(width)
}

func Border(padding, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		PaddingTop(padding / 2).
		PaddingBottom(padding / 2).
		PaddingRight(padding).
		PaddingLeft(padding).
		Width(width).
		BorderForeground(lipgloss.Color("#50C878"))
}
