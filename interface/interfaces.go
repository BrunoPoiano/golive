package interfaces

import (
	"fmt"
	"main/models"
)

func Playing(m models.MainModel) string {

	s := "Playing\n"
	s += fmt.Sprintf(" Input: %s\n", m.Input.Items[m.Input.Selected])
	s += fmt.Sprintf("Output: %s\n", m.Output.Items[m.Output.Selected])
	s += "\n Press s to Stop\n"

	return s
}

func ListItems(m models.MainModel) string {
	// The header
	s := "Select Input:\n"

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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
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
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress p to play.\n"
	s += "Press r to refresh items.\n"
	s += "Press q to quit.\n"

	// Send the UI for rendering
	return s
}
