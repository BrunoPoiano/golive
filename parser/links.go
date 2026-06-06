package parser

import (
	"fmt"
	"os/exec"
	"strings"
)

func ReturnInputList() ([]string, error) {
	var inputs []string
	pwInputs, err := exec.Command("pw-link", "-o").Output()
	if err != nil {
		return inputs, fmt.Errorf("Error getting inputs")
	}

	stdOut := string(pwInputs)

	for _, item := range strings.Split(stdOut, "\n") {
		if strings.Contains(item, "alsa_input") {
			input := strings.Split(item, ":")[0]

			if !contains(inputs, input) {
				inputs = append(inputs, input)
			}
		}
	}

	return inputs, nil
}

func ReturnOuputList() ([]string, error) {
	var inputs []string
	pwInputs, err := exec.Command("pw-link", "-i").Output()
	if err != nil {
		return inputs, fmt.Errorf("Error getting inputs")
	}

	stdOut := string(pwInputs)

	for _, item := range strings.Split(stdOut, "\n") {
		if strings.Contains(item, "alsa_output") {
			input := strings.Split(item, ":")[0]
			if !contains(inputs, input) {
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

func contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
