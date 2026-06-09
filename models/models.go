package models

import (
	"os/exec"
)

type LevelMsg string

type MainModel struct {
	Play   *exec.Cmd
	Input  Input
	Output Output
	Level  Level
	Debug  string
	Cursor int

	Padding int
	Width   int
	Height  int
}

type Level struct {
	Process *exec.Cmd
	Value   string
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
