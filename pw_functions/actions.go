package pw

import (
	"bufio"
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
	}

	if err := json.Unmarshal(pwInputs, &dump); err != nil {
		panic(err)
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

func Play(m models.MainModel) *exec.Cmd {
	input := m.Input.Items[m.Input.Selected].Info.Props.NodeName
	output := m.Output.Items[m.Output.Selected].Info.Props.NodeName

	capture := fmt.Sprintf("--capture-props=node.target=%s", input)
	playback := fmt.Sprintf("--playback-props=node.target=%s", output)

	cmd := exec.Command("pw-loopback", capture, playback)

	cmd.Start()
	return cmd
}

func MonitorChanel(p *tea.Program, input string) *exec.Cmd {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "pulse",
		"-i", input,
		"-af", "astats=metadata=1:reset=1,ametadata=print:key=lavfi.astats.Overall.Peak_level",
		"-f", "null",
		"-")

	output, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()

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

	return cmd
}

func KillProcesses(m models.MainModel) models.MainModel {
	if m.Play != nil && m.Play.Process != nil {
		m.Play.Process.Kill()
		m.Play = nil
	}
	if m.Level.Process != nil && m.Level.Process.Process != nil {
		m.Level.Process.Process.Kill()
		m.Level.Process = nil
	}

	return m
}

func GetActiveNodes(PId int) (string, string) {
	var outputId, inputId string
	pwLoopback := fmt.Sprintf("pw-loopback-%d", PId)
	cmdOutput, err := exec.Command("pw-cli", "ls", "Node").Output()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(cmdOutput), "\n")
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf("output.%s", pwLoopback)) {
			outputId = getNodeId(lines, i)
		}
		if strings.Contains(line, fmt.Sprintf("input.%s", pwLoopback)) {
			inputId = getNodeId(lines, i)
		}
	}

	return inputId, outputId
}

func getNodeId(lines []string, index int) string {
	start := max(0, index-5)

	for _, l := range lines[start : index+1] {
		if strings.Contains(l, "id ") && !strings.Contains(l, ".id") {
			before, _, found := strings.Cut(l, ",")
			if found {
				_, value, foundValue := strings.Cut(before, " ")
				if foundValue {
					return value
				}
			}
		}
	}

	return ""
}

func ChangeVolume(id int, volume float64) {
	volumeCmd := fmt.Sprintf("{ mute: false, channelVolumes: [ %f, %f ] }", volume, volume)
	exec.Command("pw-cli", "s", strconv.Itoa(id), "Props", volumeCmd).Start()
}

func RefresLists(m models.MainModel) models.MainModel {
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
