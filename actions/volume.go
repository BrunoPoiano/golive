package actions

import (
	"main/models"
	pw "main/pw_functions"
	"math"
)

var maxVolume = 2.0
var volumeRate = 0.01
var shitVolumeRate = 0.05

func Volume(action string, m models.MainModel) models.MainModel {

	switch action {
	case "m":
		m.Output.Volume.Mute = !m.Output.Volume.Mute
	case "n":
		m.Input.Volume.Mute = !m.Input.Volume.Mute
	// output volume
	case "left":
		if m.Output.Volume.Right > 0 && m.Output.Volume.Left > 0 {
			m.Output.Volume.Right = math.Max(0, m.Output.Volume.Right-volumeRate)
			m.Output.Volume.Left = math.Max(0, m.Output.Volume.Left-volumeRate)
		}

	case "shift+left":
		if m.Output.Volume.Right > 0 && m.Output.Volume.Left > 0 {
			m.Output.Volume.Right = math.Max(0, m.Output.Volume.Right-shitVolumeRate)
			m.Output.Volume.Left = math.Max(0, m.Output.Volume.Left-shitVolumeRate)
		}

	// output volume
	case "right":
		if m.Output.Volume.Right < maxVolume && m.Output.Volume.Left < maxVolume {
			m.Output.Volume.Right = math.Min(maxVolume, m.Output.Volume.Right+volumeRate)
			m.Output.Volume.Left = math.Min(maxVolume, m.Output.Volume.Left+volumeRate)
		}

	case "shift+right":
		if m.Output.Volume.Right < maxVolume && m.Output.Volume.Left < maxVolume {
			m.Output.Volume.Right = math.Min(maxVolume, m.Output.Volume.Right+shitVolumeRate)
			m.Output.Volume.Left = math.Min(maxVolume, m.Output.Volume.Left+shitVolumeRate)
		}

	// Input Volume
	case "a":
		if m.Input.Volume.Right > 0 && m.Input.Volume.Left > 0 {
			m.Input.Volume.Right = math.Max(0, m.Input.Volume.Right-volumeRate)
			m.Input.Volume.Left = math.Max(0, m.Input.Volume.Left-volumeRate)
		}

	case "A":
		if m.Input.Volume.Right > 0 && m.Input.Volume.Left > 0 {
			m.Input.Volume.Right = math.Max(0, m.Input.Volume.Right-shitVolumeRate)
			m.Input.Volume.Left = math.Max(0, m.Input.Volume.Left-shitVolumeRate)
		}

	// Input Volume
	case "d":
		if m.Input.Volume.Right < maxVolume && m.Input.Volume.Left < maxVolume {
			m.Input.Volume.Right = math.Min(maxVolume, m.Input.Volume.Right+volumeRate)
			m.Input.Volume.Left = math.Min(maxVolume, m.Input.Volume.Left+volumeRate)
		}
	case "D":
		if m.Input.Volume.Right < maxVolume && m.Input.Volume.Left < maxVolume {
			m.Input.Volume.Right = math.Min(maxVolume, m.Input.Volume.Right+shitVolumeRate)
			m.Input.Volume.Left = math.Min(maxVolume, m.Input.Volume.Left+shitVolumeRate)
		}
	}
	go pw.ChangeVolume(m, models.StreamOutput)
	return m
}
