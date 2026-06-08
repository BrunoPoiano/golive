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

func MonitorChanel(cmd *exec.Cmd, p *tea.Program, input string) {

	cmd = exec.Command(
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

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Peak_level=") {
			_, level, found := strings.Cut(line, "Peak_level=")
			if found {
				p.Send(models.LevelMsg(level))
			}
		}
	}

}
