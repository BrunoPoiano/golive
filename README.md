# go-live

A terminal-based audio routing manager for PipeWire/PulseAudio systems. Easily select audio inputs and outputs, monitor audio levels with a real-time visual meter, and create audio loopbacks from your terminal.

<img width="1038" height="321" alt="image" src="https://github.com/user-attachments/assets/638c764a-b9b5-4e23-8ae2-37f5bcf87c65" />
<img width="1037" height="357" alt="image" src="https://github.com/user-attachments/assets/a661e88d-7359-484c-a04e-6d46bd99db73" />



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

### Navigation

- **Arrow Keys / j/k**: Move cursor up and down through inputs and outputs
- **Space / Enter**: Select highlighted input or output
- **p**: Start audio loopback with selected input and output
- **s**: Stop the currently running audio loopback
- **r**: Refresh available audio devices
- **q / Ctrl+C**: Quit the application

## Dependencies

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
