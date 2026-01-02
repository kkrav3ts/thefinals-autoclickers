package main

import (
	"fmt"
	"syscall"
	"time"
)

// Version is set at build time via -ldflags
var Version = "dev"

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybd_event      = user32.NewProc("keybd_event")
)

const (
	VK_RBUTTON      = 0x02
	VK_LCONTROL     = 0xA2
	KEYEVENTF_KEYUP = 0x0002

	PingInterval = 1 * time.Second
	PollRate     = 100 * time.Millisecond
)

func isKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

func pressKey(vk int) {
	keybd_event.Call(uintptr(vk), 0, 0, 0)
	keybd_event.Call(uintptr(vk), 0, uintptr(KEYEVENTF_KEYUP), 0)
}

func main() {
	fmt.Printf("THE FINALS Auto-Ping %s\n", Version)

	var pressed bool
	var lastPing time.Time

	for {
		if isKeyPressed(VK_RBUTTON) {
			if !pressed {
				pressed = true
				lastPing = time.Now()
				pressKey(VK_LCONTROL)
			} else if time.Since(lastPing) >= PingInterval {
				pressKey(VK_LCONTROL)
				lastPing = time.Now()
			}
		} else {
			pressed = false
		}

		time.Sleep(PollRate)
	}
}
