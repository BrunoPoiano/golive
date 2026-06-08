package models

import (
	"os/exec"
)

type LevelMsg string

type MainModel struct {
	PlayProcess  *exec.Cmd
	LevelProcess *exec.Cmd
	Input        Input
	Output       Output
	Debug        string
	Level        string
	Cursor       int
}

type Input struct {
	Items    []string
	Selected int
}

type Output struct {
	Items    []string
	Selected int
}

type PwLinks string

const (
	OutputList PwLinks = "-i"
	InputList  PwLinks = "-o"
)
