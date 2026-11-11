# THE FINALS Auto-Clicker Tools

Level up your gameplay with smart automation. Three lightweight tools designed for THE FINALS that give you a competitive edge — without touching game files.

---

## 🎯 Auto-Ping

**Never miss a callout again.** Hold right-click to aim, and enemies get pinged automatically every second. Perfect for solo queue where communication wins games.

![Auto-Ping Demo](https://github.com/user-attachments/assets/bff17d8c-ceb6-458f-b3ba-b447ce5a8e2a)

### Why use it?

- 🏆 **Win more fights** — Your team always knows where enemies are
- 🎮 **Zero effort** — Just aim like you normally would
- 🤝 **Solo queue friendly** — Communicate without a mic

---

## 🔫 Auto-Shot

**Maximize your fire rate.** Transforms any semi-auto or burst weapon into a consistent clicking machine with human-like timing. Hold left mouse button and let the tool do the clicking.

![Auto-Shot Demo](https://github.com/user-attachments/assets/9637f8b4-4a63-41e1-90ed-ceec852d2446)

### Why use it?

- ⚡ **Faster firing** — Consistent clicks without finger fatigue
- 🎯 **Better accuracy** — Focus on aiming, not clicking
- 🕹️ **Natural timing** — Randomized delays that feel human

---

## ⏱️ Delay Checker

**Measure your click timing.** A simple diagnostic tool that displays the duration of your left mouse button clicks in milliseconds. Useful for testing your mouse, analyzing your click patterns, or calibrating the Auto-Shot tool.

![Delay checker](https://github.com/user-attachments/assets/7dda3dd5-565e-497a-91f8-6a000baf89da)

### Why use it?

- 🔬 **Debug your setup** — Verify your mouse is registering clicks correctly
- 📊 **Analyze patterns** — See how long you hold clicks on average
- 🛠️ **Calibrate tools** — Fine-tune Auto-Shot timing based on your click behavior

---

## 🎮 Download & Install

1. Go to the [**Releases**](../../releases) page
2. Download the `.exe` file for the tool you want:
    - `thefinals-autoping-vX.X.X.exe` — for auto-pinging
    - `thefinals-autoshot-vX.X.X.exe` — for auto-shooting
    - `thefinals-delaychecker-vX.X.X.exe` — for click timing analysis
3. Run it — no installation needed!

> **Note:** Windows only. Fully open-source, no malicious code!

> ⚠️ **Windows Defender Warning:** These executables are not signed with a Microsoft certificate, so Windows Defender or SmartScreen may flag them as unrecognized. This is normal for open-source tools. You can safely allow the app to run — the source code is fully available for review.

---

## 📖 How to Use

### Auto-Ping Setup

1. **Run** `thefinals-autoping.exe`
2. **Press** the keyboard key you use for pinging in-game (e.g., `Z` or `X`)
3. **Play!** Hold right-click (aim) and enemies get pinged automatically

Optional flags tune the rhythm without rebuilding. Defaults match the previous release: a ping about once a second, plus 70–100ms of jitter.

```bash
thefinals-autoping.exe -interval 1s -hold 10ms -mean 80 -stddev 5 -min 70 -max 100
```

| Action  | Key                      |
|---------|--------------------------|
| Trigger | Right Mouse Button (Aim) |
| Ping    | Your selected key        |

---

### Auto-Shot Setup

This tool requires a quick one-time setup in your game settings:

#### Step 1: Change your in-game shoot key

1. Open THE FINALS → **Settings** → **Keybinds**
2. Change **"Fire"** from `Left Mouse Button` to any keyboard key (e.g., `L`)
3. Apply and close settings

#### Step 2: Configure the tool

1. **Run** `thefinals-autoshot.exe`
2. **Press** the same key you just mapped for shooting (e.g., `L`)
3. The tool will confirm your selection

#### Step 3: Play!

- **Hold Left Mouse Button** to shoot — the tool rapidly presses your mapped key
- Release to stop shooting

Each click holds the key for one random sample, then waits for a second sample. Defaults stay in the 50–80ms band. Raise `-min` and `-max` if a patch gives a weapon a slower fire rate, or lower them if your clicks are not keeping up.

```bash
thefinals-autoshot.exe -mean 60 -stddev 5 -min 50 -max 80
```

| Action  | Key                      |
|---------|--------------------------|
| Trigger | Left Mouse Button (Hold) |
| Shoot   | Your selected key        |

> **Tip:** This works great with semi-auto weapons like pistols, revolvers, or burst weapons!

---

### Delay Checker Setup

1. **Run** `thefinals-delaychecker.exe`
2. **Click** the left mouse button anywhere
3. **See** the click duration displayed in milliseconds

```
Monitoring left mouse button clicks (high-precision mode)...
Press and release the left mouse button to see delays.

[Click #1] First press detected
[Click #1] Hold duration: 87.40 ms (avg 87.40 ms)

[Click #2] Interval since last press: 190.20 ms (avg 190.20 ms)
[Click #2] Hold duration: 95.10 ms (avg 91.25 ms)
```

> **Tip:** Use this to understand your natural click timing before configuring Auto-Shot!

---

## ✅ Safe to Use

All tools are **external** and work by monitoring keyboard state and synthesizing keystrokes with `SendInput`. They:

- ❌ Do NOT read game memory
- ❌ Do NOT modify game files
- ❌ Do NOT inject into the game process
- ✅ Work like a hardware macro keyboard

---

## 🛠️ Build from Source (Advanced)

For developers who want to compile from source:

**Requirements:** [Go 1.26+](https://golang.org/dl/)

```bash
# Build Auto-Ping
go build -o thefinals-autoping.exe ./cmd/autoping

# Build Auto-Shot
go build -o thefinals-autoshot.exe ./cmd/autoshot

# Build Delay Checker
go build -o thefinals-delaychecker.exe ./cmd/delayChecker
```

See [Windows Virtual-Key Codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes) for key code reference.

---

## 💬 Feedback & Support

Have an idea, found a bug, or need help? [Open an issue](../../issues) — I'd love to hear from you!

---
<p align="center">
  <strong>Maintained with ❤️ by kkrav3ts. Good luck, contestant!
</p>
