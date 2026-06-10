package pw

import (
	"bufio"
	"fmt"
	"main/models"
	"os/exec"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func ReturnList(t models.PwLinks) ([]string, error) {

	var inputs []string
	pwInputs, err := exec.Command("pw-link", string(t)).Output()
	if err != nil {
		return inputs, fmt.Errorf("Error getting inputs")
	}

	stdOut := string(pwInputs)
	pwtype := "alsa_input"

	if string(t) == string(models.OutputList) {
		pwtype = "alsa_output"
	}

	for item := range strings.SplitSeq(stdOut, "\n") {
		if strings.Contains(item, pwtype) {
			input, _, found := strings.Cut(item, ":")

			if found && !slices.Contains(inputs, input) {
				inputs = append(inputs, input)
			}
		}
	}

	return inputs, nil
}

func Play(input, output string) *exec.Cmd {
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

func ChangeVolume(volume models.Volume) {
	volumeCmd := fmt.Sprintf("{ mute: false, channelVolumes: [ %f, %f ] }", volume.Value, volume.Value)
	exec.Command("pw-cli", "s", volume.NodeId, "Props", volumeCmd).Start()
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
