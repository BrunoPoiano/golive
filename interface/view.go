package interfaces

import (
	"main/models"
	"strings"
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

	right.WriteString("Increase Input  Vol: d\nDecrease Input  Vol: a\nIncrease Output Vol: right\nDecrease Output Vol: left")

	if m.Play.Cmd != nil {
		left.WriteString(Playing(m))
		right.WriteString("\nx: Stop\nq: quit")
	} else {
		left.WriteString(ListItems(m))
		right.WriteString("\np: play\nr: refresh lists\nq: quit")
	}

	return header, left.String(), right.String()

}
