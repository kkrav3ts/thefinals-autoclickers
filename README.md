# THE FINALS Auto-Ping

A lightweight Windows utility that automatically pings enemies while you're aiming in THE FINALS. It greatly improves communication within the team and helps you win more games even in solo queue.

## Demo

![Demo](https://github.com/user-attachments/assets/bff17d8c-ceb6-458f-b3ba-b447ce5a8e2a)

## How It Works

When you hold the **right mouse button** (aim down sights), the tool automatically presses the **Middle Mouse Button** (ping) for you:

1. **Initial press** — Pings immediately when you start aiming
2. **Repeat** — Continues to ping every 1 second while you hold right-click
3. **Release** — Stops when you release the right mouse button

This keeps enemies marked without interrupting your aim.

## Download

Download the latest release from the [Releases](../../releases) page.
Only Windows supported.

### Build from Source

- Requires Go 1.25. Install Go from [golang.org](https://golang.org/dl/).
- Clone this repository.
- Run the following command in the terminal:

```bash
# On Windows
go build -o autoping.exe main.go

# Cross-compile from macOS/Linux
GOOS=windows GOARCH=amd64 go build -o autoping.exe main.go
```

## Usage

1. Run `autoping.exe`
2. Launch THE FINALS
3. Aim at enemies — they get pinged automatically!
4. Press `Ctrl+C` in the terminal to stop

## Configuration

Edit the constants in `main.go` to customize behavior:

```go
PingInterval   = 1 * time.Second // Time between pings
PollRateActive = 50 * time.Millisecond // Polling speed while aiming
PollRateIdle = 200 * time.Millisecond // Polling speed when idle
KeyPressDelay = 10 * time.Millisecond // Delay between key press/release
```

## Key Bindings

| Action  | Default Key         |
|---------|---------------------|
| Trigger | Right Mouse Button  |
| Ping    | Middle Mouse Button |

To change keys, modify `VK_RBUTTON` and `VK_MBUTTON` in `main.go` using [Windows Virtual-Key Codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes).

## License

MIT
