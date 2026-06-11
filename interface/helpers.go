package interfaces

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

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
	value = math.Max(0, math.Min(66, value))

	live := strings.Repeat("|", int(value))

	return fmt.Sprintf("%s\n%s\n%s", ruler, live, ruler)
}
