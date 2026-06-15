package interfaces

import (
	"fmt"
	"main/models"
	"math"
	"strconv"
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

	RMSLevel, err := strconv.ParseFloat(m.Level.RMSLevel, 64)
	if math.IsInf(RMSLevel, 0) {
		fmt.Fprintf(&s, "Raw RMS Level: 0.0 dBFS\n")
	} else if err == nil {
		fmt.Fprintf(&s, "Raw RMS Level: %.1f dBFS\n", RMSLevel)
	}

	PeakLevel, err := strconv.ParseFloat(m.Level.PeakLevel, 64)
	if math.IsInf(RMSLevel, 0) {
		fmt.Fprintf(&s, "Raw Signal Peak: 0.0 dBFS")
		fmt.Fprintf(&s, "\n%s", generateMeter(0.0))
	} else if err == nil {
		fmt.Fprintf(&s, "Raw Signal Peak: %.1f dBFS", PeakLevel)
		fmt.Fprintf(&s, "\n%s", generateMeter(PeakLevel))
	}

	fmt.Fprintf(&s, "\n\n Input: %s", (m.Input.Items[m.Input.Selected].Info.Props.NodeDescription))
	fmt.Fprintf(&s, "\nOutput: %s", (m.Output.Items[m.Output.Selected].Info.Props.NodeDescription))

	//gain
	fmt.Fprintf(&s, "\n\n Input Gain: (%.1fdb)\n", AmplitudeToDB(m.Input.Volume.Value))
	s.WriteString(generateGain(m.Input.Volume.Value, 200))
	fmt.Fprintf(&s, "\nOutput Gain: (%.1fdb)\n", AmplitudeToDB(m.Output.Volume.Value))
	s.WriteString(generateGain(m.Output.Volume.Value, 200))

	return s.String()
}

func ListItems(m models.MainModel) string {

	inputPercent, _ := generateVolume((m.Input.Volume.Value))
	outputPercent, _ := generateVolume((m.Output.Volume.Value))
	var s strings.Builder
	// s = fmt.Sprintf("debug: %s\n", m.Debug)
	fmt.Fprintf(&s, "Inputs%s\n", inputPercent)
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

	fmt.Fprintf(&s, "\nOutputs%s\n", outputPercent)

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
