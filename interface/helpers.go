package interfaces

import (
	"fmt"
	"math"
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

func generateMeter(peakLevel float64) string {

	value := peakLevel
	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = 1
	}

	value = math.Floor(value * -1)
	value = math.Max(0, math.Min(66, value))

	live := strings.Repeat("|", int(value))

	return fmt.Sprintf("%s\n%s\n%s", ruler, live, ruler)
}

func generateGain(volume, maxValue int) string {

	var live string
	var percent string
	live += strings.Repeat("█", volume/4)
	live += strings.Repeat("░", ((maxValue / 4) - (volume / 4)))

	if volume > 150 {
		percent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4800")).
			Render(fmt.Sprintf(" (%d%%)", volume))
	} else if volume > 100 {
		percent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FF00")).
			Render(fmt.Sprintf(" (%d%%)", volume))
	} else {
		percent = fmt.Sprintf(" (%d%%)", volume)
	}
	live += percent
	return fmt.Sprintf("%s", live)
}

func calcVolumePercent(v float64) int {
	return int(v * 100 / 1)
}

func AmplitudeToDB(v float64) float64 {

	calc := 20 * math.Log10(v)
	if math.IsInf(calc, 0) {
		return 20 * math.Log10(0.01)
	}

	return calc
}
