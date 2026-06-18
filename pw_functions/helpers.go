package pw

import (
	"main/models"
	"os/exec"
	"strconv"
	"strings"
)

func stop(p *models.Action) error {
	if p == nil {
		return nil
	}

	if p.Cancel != nil {
		p.Cancel()
	}

	if p.Cmd == nil {
		return nil
	}

	err := p.Cmd.Wait()

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			p.Cmd = nil
			p.Cancel = nil
			return nil
		}

		return err
	}

	p.Cmd = nil
	p.Cancel = nil

	return nil
}

func parseLevels(line, search string, defaultValue float64) float64 {
	if strings.Contains(line, search) {
		_, value, found := strings.Cut(line, search)
		if found {
			value64, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return defaultValue
			}
			return value64
		}
	}

	return defaultValue
}
