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

func Play(m models.MainModel) models.MainModel {

	input := m.Input.Items[m.Input.Selected].Info.Props.NodeName
	output := m.Output.Items[m.Output.Selected].Info.Props.NodeName

	capture := fmt.Sprintf("--capture-props=node.target=%s", input)
	playback := fmt.Sprintf("--playback-props=node.target=%s", output)

	ctx, cancel := context.WithCancel(context.Background())

	m.Play.Cancel = cancel
	m.Play.Cmd = exec.CommandContext(ctx, "pw-loopback", capture, playback)

	m.Play.Cmd.Start()
	go ChangeVolume(m.Input.Items[m.Input.Selected].Id, m.Input.Volume)
	go ChangeVolume(m.Output.Items[m.Output.Selected].Id, m.Output.Volume)
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
		panic(err)
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

func KillProcesses(m models.MainModel) models.MainModel {
	if err := stop(&m.Play); err != nil {
		fmt.Println(err)
	}

	if err := stop(&m.Level.Action); err != nil {
		fmt.Println(err)
	}

	return m
}

func ChangeVolume(id int, volume float64) error {
	volumeCmd := fmt.Sprintf("{ mute: false, channelVolumes: [ %f, %f ] }", volume, volume)
	return exec.Command("pw-cli", "s", strconv.Itoa(id), "Props", volumeCmd).Start()
}

func RefreshLists(m models.MainModel) models.MainModel {
	inputsList, err := ReturnList(models.InputList)
	if err != nil {
		panic(err.Error())
	}
	outputList, err := ReturnList(models.OutputList)
	if err != nil {
		panic(err.Error())
	}

	m.Input.Items = inputsList
	m.Output.Items = outputList

	return m
}
