package parser

import (
	"fmt"
	"main/models"
	"os/exec"
	"slices"
	"strings"
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
