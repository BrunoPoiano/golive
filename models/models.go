package models

import (
	"context"
	"os/exec"
)

type LevelMsg struct {
	PeakLevel string
	RMSLevel  string
}
type ErrorMsg error

type MainModel struct {
	Play   Action
	Input  Input
	Output Output
	Level  Level
	Debug  string
	Error  error
	Cursor int

	Padding int
	Width   int
	Height  int
}

type Action = struct {
	Cmd    *exec.Cmd
	Cancel context.CancelFunc
}

type PwDump struct {
	Id   int `json:"id"`
	Info struct {
		Props struct {
			ObjectId        int    `json:"node.object.id"`
			NodeName        string `json:"node.name"`
			NodeDescription string `json:"node.description"`
			Stream          string `json:"api.alsa.pcm.stream"`
		} `json:"props"`
	} `json:"info"`
}

type Volume struct {
	Value float64
	Mute  bool
}
type Level struct {
	Action    Action
	PeakLevel string
	RMSLevel  string
}

type Input struct {
	Items    []PwDump
	Selected int
	Volume   Volume
}

type Output struct {
	Items    []PwDump
	Selected int
	Volume   Volume
}

type PwLinks string

const (
	PlaybackList PwLinks = "playback"
	CaptureList  PwLinks = "capture"
)

type AppColor string

const (
	Danger    AppColor = "#FF4800"
	Attention AppColor = "#F1FF00"
	Success   AppColor = "#50C878"
)
