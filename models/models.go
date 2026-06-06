package models

import "os/exec"

type MainModel struct {
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
