package interfaces

import (
	"fmt"
	"main/models"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

var ruler string = fmt.Sprintf("%s%s%s", lipgloss.NewStyle().
	Foreground(lipgloss.Color(string(models.Danger))).
	Render(fmt.Sprintf("%s6", strings.Repeat("-", 5))),
	lipgloss.NewStyle().
		Foreground(lipgloss.Color(string(models.Attention))).
		Render(fmt.Sprintf("%s12%s", strings.Repeat("-", 5), strings.Repeat("-", 4))),
	lipgloss.NewStyle().
		Foreground(lipgloss.Color(string(models.Success))).
		Render(fmt.Sprintf("18%s24%s30%s36%s42%s48%s54%s60%s",
			strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
			strings.Repeat("-", 4), strings.Repeat("-", 4), strings.Repeat("-", 4),
			strings.Repeat("-", 6), strings.Repeat("-", 6))),
)

var dangerBar = lipgloss.NewStyle().
	Foreground(lipgloss.Color(string(models.Danger))).
	Render("|")

func generateMeter(peakLevel, max float64) string {

	value := peakLevel
	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = 1
	}

	value = math.Floor(value * -1)
	value = math.Max(0, math.Min(66, value))

	live := strings.Repeat("|", int(value))

	intMax := int(math.Floor(max * -1))

	if intMax > 0 && len(live) > intMax {
		live = live[:intMax-1] + dangerBar + live[(intMax+1):]
	}

	return fmt.Sprintf("%s\n%s\n%s", ruler, live, ruler)
}

func generateGain(volume float64, maxValue int) string {
	percent, live := generateVolume(volume)
	volumePercent := calcVolumePercent(volume)

	live += strings.Repeat("░", ((maxValue / 4) - (volumePercent / 4)))
	live += percent
	return fmt.Sprintf("%s", live)
}

func generateVolume(v float64) (string, string) {
	volume := calcVolumePercent(v)

	switch {
	case volume > 150:
		return lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Danger))).
				Render(fmt.Sprintf(" [%d%%]", volume)),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Danger))).
				Render(strings.Repeat("█", volume/4))

	case volume > 100:
		return lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Attention))).
				Render(fmt.Sprintf(" [%d%%]", volume)),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Attention))).
				Render(strings.Repeat("█", volume/4))
	default:
		return lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Success))).
				Render(fmt.Sprintf(" [%d%%]", volume)),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(string(models.Success))).
				Render(strings.Repeat("█", volume/4))
	}

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
