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
	fmt.Fprintf(&s, "%s\n", generateMeter(m.Level.Value))
	fmt.Fprintf(&s, "debug %d\n", m.Input.Items[m.Input.Selected].Id)
	s.WriteString("\nPlaying\n")
	fmt.Fprintf(&s, " Input: %d%% %s\n", int(m.Input.Volume*100/1), (m.Input.Items[m.Input.Selected].Info.Props.NodeDescription))
	fmt.Fprintf(&s, "Output: %d%% %s\n", int(m.Output.Volume*100/1), (m.Output.Items[m.Output.Selected].Info.Props.NodeDescription))

	return s.String()
}

func ListItems(m models.MainModel) string {
	// The header

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

func fixName(name string) string {
	splited := strings.Split(name, ".")

	audioType := strings.ReplaceAll(splited[len(splited)-1], "-", " ")
	return fmt.Sprintf("%13s | %s ", audioType, splited[1])
}

var ruler string = fmt.Sprintf("%s%s%s", lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF4800")).
	Render(fmt.Sprintf("%s6", strings.Repeat("-", 5))),
	lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F1FF00")).
		Render(fmt.Sprintf("%s12%s", strings.Repeat("-", 5), strings.Repeat("-", 4))),
	lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50C878")).
		Render(fmt.Sprintf("18%s24%s30%s36%s42%s48%s54%s60%s",
			strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
			strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
			strings.Repeat("-", 6), strings.Repeat("-", 6))),
)

func generateMeter(peakLevel string) string {

	value, err := strconv.ParseFloat(peakLevel, 32)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		value = 1
	}

	value = math.Floor(value * -1)
	if value < 1 {
		value = 1
	}
	if value > 66 {
		value = 66
	}

	live := strings.Repeat("|", int(value))

	return fmt.Sprintf("%s\n%s\n%s", ruler, live, ruler)
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
