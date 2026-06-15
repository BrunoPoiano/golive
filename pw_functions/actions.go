package pw

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"main/models"
	"os/exec"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func ReturnList(stream models.PwLinks) ([]models.PwDump, error) {

	var list []models.PwDump
	var dump []models.PwDump

	pwInputs, err := exec.Command("pw-dump").Output()
	if err != nil {
		return list, fmt.Errorf("Error getting list")
	}

	if err := json.Unmarshal(pwInputs, &dump); err != nil {
		return list, fmt.Errorf("Error parsing list")
	}

	for _, item := range dump {
		if item.Info.Props.Stream == string(stream) {
			list = append(list, item)
		}
	}

	return list, nil
}

func Play(p *tea.Program, m models.MainModel) models.MainModel {

	input := m.Input.Items[m.Input.Selected].Info.Props.NodeName
	output := m.Output.Items[m.Output.Selected].Info.Props.NodeName

	capture := fmt.Sprintf("--capture-props=node.target=%s", input)
	playback := fmt.Sprintf("--playback-props=node.target=%s", output)

	ctx, cancel := context.WithCancel(context.Background())

	m.Play.Cancel = cancel
	m.Play.Cmd = exec.CommandContext(ctx, "pw-loopback", capture, playback)

	err := m.Play.Cmd.Start()
	if err != nil {
		p.Send(models.ErrorMsg(err))
		return m
	}

	go func() {
		err := ChangeVolume(m.Input.Items[m.Input.Selected].Id, m.Input.Volume)
		if err != nil {
			p.Send(models.ErrorMsg(err))
		}
		err = ChangeVolume(m.Output.Items[m.Output.Selected].Id, m.Output.Volume)
		if err != nil {
			p.Send(models.ErrorMsg(err))
		}

	}()

	return m
}

func MonitorChannel(p *tea.Program, m models.MainModel) models.MainModel {

	input := m.Input.Items[m.Input.Selected].Info.Props.NodeName
	ctx, cancel := context.WithCancel(context.Background())

	m.Level.Action.Cancel = cancel
	m.Level.Action.Cmd = exec.CommandContext(ctx,
		"ffmpeg",
		"-f", "pulse",
		"-i", input,
		"-af", "astats=metadata=1:reset=1,ametadata=mode=print",
		"-f", "null",
		"-")

	output, err := m.Level.Action.Cmd.StderrPipe()
	if err != nil {
		p.Send(models.ErrorMsg(fmt.Errorf("Error getting level")))
	}

	m.Level.Action.Cmd.Start()
	scanner := bufio.NewScanner(output)

	go func() {
		level := models.LevelMsg{
			PeakLevel: "0.0",
			RMSLevel:  "0.0",
		}
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "RMS_level=") {
				_, rmsLevel, found := strings.Cut(line, "RMS_level=")
				if found {
					level.RMSLevel = rmsLevel
				}
			}
			if strings.Contains(line, "Peak_level=") {
				_, peakLevel, found := strings.Cut(line, "Peak_level=")
				if found {
					level.PeakLevel = peakLevel
				}
			}
			p.Send(models.LevelMsg(level))
		}
	}()

	if err := scanner.Err(); err != nil {
		p.Send(models.ErrorMsg(fmt.Errorf("Error getting level")))
	}

	return m
}

func KillProcesses(p *tea.Program, m models.MainModel) models.MainModel {
	if err := stop(&m.Play); err != nil {
		p.Send(models.ErrorMsg(fmt.Errorf("Error killing play process")))
	}

	if err := stop(&m.Level.Action); err != nil {
		p.Send(models.ErrorMsg(fmt.Errorf("Error level meter process")))
	}

	return m
}

func ChangeVolume(id int, volume models.Volume) error {

	mute := "false"
	if volume.Mute {
		mute = "true"
	}

	volumeCmd := fmt.Sprintf("{ mute: %s, channelVolumes: [ %f, %f ] }", mute, volume.Value, volume.Value)
	start := exec.Command("pw-cli", "s", strconv.Itoa(id), "Props", volumeCmd)
	err := start.Start()
	if err != nil {
		return err
	}
	return start.Wait()
}

func RefreshLists(m models.MainModel) (models.MainModel, error) {
	var err error
	m.Input.Selected = 0
	m.Output.Selected = 0
	m.Input.Items, err = ReturnList(models.CaptureList)
	if err != nil {
		return m, fmt.Errorf("Error getting inputs")
	}
	m.Output.Items, err = ReturnList(models.PlaybackList)
	if err != nil {
		return m, fmt.Errorf("Error getting outputs")
	}
	return m, nil
}
