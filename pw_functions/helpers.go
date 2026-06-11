package pw

import (
	"main/models"
	"os/exec"
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
