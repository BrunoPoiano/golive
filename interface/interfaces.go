package interfaces

import (
	"fmt"
	"main/models"
	"math"
	"strconv"
	"strings"
)

func Playing(m models.MainModel) string {

	var s string
	s += fmt.Sprintf("Meter: %s\n", generateMeter(m.Level))
	s += "\nPlaying\n"
	s += fmt.Sprintf(" Input: %s\n", fixName(m.Input.Items[m.Input.Selected]))
	s += fmt.Sprintf("Output: %s\n", fixName(m.Output.Items[m.Output.Selected]))
	s += "\n Press s to Stop\n"

	return s
}

func ListItems(m models.MainModel) string {
	// The header

	var s string
	// s = fmt.Sprintf("debug: %s\n", m.Debug)
	s += "Select Input:\n"

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, fixName(choice))
	}

	s += "\nSelect output:\n"

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, fixName(choice))
	}

	// The footer
	s += "\nPress p to play.\n"
	s += "Press r to refresh items.\n"
	s += "Press q to quit.\n"

	// Send the UI for rendering
	return s
}

func fixName(name string) string {
	splited := strings.Split(name, ".")

	return fmt.Sprintf("%s | %s ", splited[len(splited)-1], splited[1])
}

func generateMeter(peakLevel string) string {

	value, err := strconv.ParseFloat(peakLevel, 32)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		value = 1
	}

	value = math.Floor(value * -1)
	if value < 1 {
		value = 1
	}
	if value > 80 {
		value = 80
	}

	var meter string

	ruler := fmt.Sprintf("%s6%s12%s18%s24%s30%s36%s42%s48%s54%s",
		strings.Repeat("-", 5), strings.Repeat("-", 5), strings.Repeat("-", 4),
		strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
		strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
		strings.Repeat("-", 6))

	live := strings.Repeat("|", int(value))

	meter += fmt.Sprintf("\n%s\n%s\n%s", ruler, live, ruler)
	return meter
}
