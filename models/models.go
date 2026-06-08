package models

import "os/exec"

type MainModel struct {
	Debug   string
	Input   Input
	Output  Output
	Process *exec.Cmd
	Cursor  int
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
