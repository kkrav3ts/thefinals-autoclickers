# THE FINALS Auto-Clicker

Automatically ping enemies while aiming in THE FINALS. Improve team communication and win more games — even in solo queue.

![Ping](https://github.com/user-attachments/assets/bff17d8c-ceb6-458f-b3ba-b447ce5a8e2a)
![Shoot](https://github.com/user-attachments/assets/9637f8b4-4a63-41e1-90ed-ceec852d2446)

## ✨ Features

- **Auto-ping while aiming** — Hold right-click, enemies get pinged automatically
- **Repeating pings** — Keeps pinging every second while you aim
- **Lightweight** — Minimal CPU usage, runs quietly in the background
- **Safe** — Does not modify game files or interact with the game process

## 🎮 Download

1. Go to the [Releases](../../releases) page
2. Download the latest `.exe` file
3. Run it — no installation needed

> **Note:** Windows only.

## 📖 Usage

1. Run the executable
2. When prompted, press the keyboard(!) key you want to use for pinging (this should match your in-game ping key)
3. The tool will confirm which key was selected
4. Hold right-click (aim) in-game to automatically ping enemies every second
5. Press `Ctrl+C` or close the window to exit

## ⚙️ Controls

| Action  | Key                |
|---------|--------------------|
| Trigger | Right Mouse Button |
| Ping    | User Selected      |

> **Note:** When you first run the tool, it will prompt you to press the key you want to use for pinging. Make sure this matches your in-game ping key binding.

## 🛠️ Advanced: Build from Source

For developers who want to compile from source:

**Requirements:** [Go 1.25+](https://golang.org/dl/)

```bash
# Windows
go build -o autoping.exe ./cmd/autoping

# macOS / Linux (cross-compile)
GOOS=windows GOARCH=amd64 go build -o autoping.exe ./cmd/autoping
```

**Customization:** 
- To change the aim trigger key, modify the `aimKey` variable in `cmd/autoping/main.go`. Default is `0x02` (Right Mouse Button).
- The ping key is selected at runtime, but you can modify the `KeyNames` map in `internal/keyboard/keyboard.go` to add/remove supported keys.
- Timing intervals can be adjusted in the `main()` function in `cmd/autoping/main.go`.

See [Windows Virtual-Key Codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes) for reference.

## 💬 Feedback & Suggestions

Have an idea or found a bug? Feel free to [open an issue](../../issues) or contact me directly.
