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

type PwDump struct {
	Id   int `json:"id"`
	Info struct {
		Props struct {
			NodeId          int    `json:"node.driver-id"`
			NodeName        string `json:"node.name"`
			NodeNick        string `json:"node.nick"`
			NodeDescription string `json:"node.description"`
			NodeGroup       string `json:"node.group"`
		} `json:"props"`
	} `json:"info"`
}

type Level struct {
	Process *exec.Cmd
	Value   string
}

type Input struct {
	Items    []PwDump
	Selected int
	Volume   float64
}

type Output struct {
	Items    []PwDump
	Selected int
	Volume   float64
}

type PwLinks string

const (
	OutputList PwLinks = "-i"
	InputList  PwLinks = "-o"
)
