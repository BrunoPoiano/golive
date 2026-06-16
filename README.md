# go-live

A terminal-based audio routing manager for PipeWire/PulseAudio systems. Easily select audio inputs and outputs, monitor audio levels with a real-time visual meter, and create audio loopbacks from your terminal.

<img width="1001" height="327" alt="image" src="https://github.com/user-attachments/assets/5d288441-1561-48e7-9ea0-111e9e90975f" />
<img width="1001" height="364" alt="image" src="https://github.com/user-attachments/assets/2a87f337-7167-4acd-8662-89c0a154c66c" />

## Features

- **Audio Loopback**: Create loopback connections between selected input and output devices
- **Real-time Level Monitoring**: Visual meter showing peak audio levels
- **Terminal UI**: Clean, intuitive interface built with Bubble Tea

## Requirements

- Go 1.26.3 or later
- PipeWire/PulseAudio
- ffmpeg

## Installation

1. Clone the repository:

```bash
git clone https://github.com/BrunoPoiano/golive.git
cd golive
```

2. Install dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o golive
```

4. Run the application:

```bash
./golive
```

## Usage

### Keybindings

#### Navigation

- **Arrow Keys / j/k**: Move cursor up and down through devices
- **Space / Enter**: Select highlighted input or output

#### Audio Control

- **p**: Start audio loopback with selected input and output
- **x**: Stop the currently running audio loopback
- **a**: Decrease input volume (gain)
- **d**: Increase input volume (gain)
- **← (Left Arrow)**: Decrease output volume (gain)
- **→ (Right Arrow)**: Increase output volume (gain)
- **n**: Toggle input mute
- **m**: Toggle output mute

#### Utilities

- **r**: Refresh available audio devices
- **q / Ctrl+C**: Quit the application

## Dependencies

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss v2](https://github.com/charmbracelet/lipgloss) - Terminal styling and layout
- ffmpeg
- Standard Go libraries (context, os, exec, bufio, fmt, strings, math, encoding/json)

## License

MIT License - Copyright (c) 2026 Bruno Poiano

See [LICENSE](LICENSE) file for details.
