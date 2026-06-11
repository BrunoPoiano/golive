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

func ReturnList(t models.PwLinks) ([]models.PwDump, error) {

	var list []models.PwDump
	var dump []models.PwDump

	pwInputs, err := exec.Command("pw-dump").Output()
	if err != nil {
		return list, fmt.Errorf("Error getting list")
	}

	if err := json.Unmarshal(pwInputs, &dump); err != nil {
		return list, fmt.Errorf("Error parsing list")
	}

	pwtype := "alsa_input"

	if string(t) == string(models.OutputList) {
		pwtype = "alsa_output"
	}

	for _, item := range dump {
		if strings.Contains(item.Info.Props.NodeName, pwtype) {
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
		p.Send(models.ErrorMsg("Error starting player"))
		return m
	}

	go func() {
		err := ChangeVolume(m.Input.Items[m.Input.Selected].Id, m.Input.Volume)
		if err != nil {
			p.Send(models.ErrorMsg("Error changing input volume"))
		}
		ChangeVolume(m.Output.Items[m.Output.Selected].Id, m.Output.Volume)
		if err != nil {
			p.Send(models.ErrorMsg("Error changing output volume"))
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
		"-af", "astats=metadata=1:reset=1,ametadata=print:key=lavfi.astats.Overall.Peak_level",
		"-f", "null",
		"-")

	output, err := m.Level.Action.Cmd.StderrPipe()
	if err != nil {
		p.Send(models.LevelMsg("Error getting level"))
	}

	m.Level.Action.Cmd.Start()
	scanner := bufio.NewScanner(output)

	if err := scanner.Err(); err != nil {
		p.Send(models.LevelMsg("Error getting level"))
	}

	go func() {
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "Peak_level=") {
				_, level, found := strings.Cut(line, "Peak_level=")
				if found {
					p.Send(models.LevelMsg(level))
				}
			}
		}
	}()

	return m
}

func KillProcesses(p *tea.Program, m models.MainModel) models.MainModel {
	if err := stop(&m.Play); err != nil {
		p.Send(models.ErrorMsg("Error killing play process"))
	}

	if err := stop(&m.Level.Action); err != nil {
		p.Send(models.ErrorMsg("Error level meter process"))
	}

	return m
}

func ChangeVolume(id int, volume float64) error {
	volumeCmd := fmt.Sprintf("{ mute: false, channelVolumes: [ %f, %f ] }", volume, volume)
	return exec.Command("pw-cli", "s", strconv.Itoa(id), "Props", volumeCmd).Start()
}

func RefreshLists(m models.MainModel) (models.MainModel, error) {
	var err error
	m.Input.Selected = 0
	m.Output.Selected = 0
	m.Input.Items, err = ReturnList(models.InputList)
	if err != nil {
		return m, fmt.Errorf("Error getting inputs")
	}
	m.Output.Items, err = ReturnList(models.OutputList)
	if err != nil {
		return m, fmt.Errorf("Error getting outputs")
	}
	return m, nil
}
