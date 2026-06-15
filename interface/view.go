package interfaces

import (
	"main/models"
	"strings"

	"charm.land/lipgloss/v2"
)

func View(m models.MainModel) (string, string, string) {

	header := Header()
	var left strings.Builder
	var right strings.Builder

	if m.Error != nil {
		left.WriteString(m.Error.Error())
		right.WriteString("r: retry")

		return header, left.String(), right.String()
	}

	if len(m.Input.Items) == 0 || len(m.Output.Items) == 0 {

		if len(m.Input.Items) == 0 {
			left.WriteString("No Inputs found\n")
		}
		if len(m.Output.Items) == 0 {
			left.WriteString("No Outputs found")
		}

		right.WriteString("r: retry")

		return header, left.String(), right.String()
	}

	right.WriteString(actions(m))

	if m.Play.Cmd != nil {
		left.WriteString(Playing(m))
		right.WriteString("\nx: Stop\nq: quit")
	} else {
		left.WriteString(ListItems(m))
		right.WriteString("\np: play\nr: refresh lists\nq: quit")
	}

	return header, left.String(), right.String()

}

func actions(m models.MainModel) string {

	var input strings.Builder
	var output strings.Builder

	input.WriteString("Input\na d: Volume")

	if m.Input.Volume.Mute {
		input.WriteString("\nUnmute: n")
	} else {
		input.WriteString("\nMute: n")
	}

	output.WriteString("Output\n← →: Volume")
	if m.Output.Volume.Mute {
		output.WriteString("\nUnmute: m\n")
	} else {
		output.WriteString("\nMute: m\n")
	}

	left := input.String()
	if m.Width > 110 {
		left = lipgloss.NewStyle().MarginRight(2).Render(input.String())
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, left, output.String())

}
