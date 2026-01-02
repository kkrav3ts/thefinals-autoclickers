// Package main provides an auto-ping utility for THE FINALS.
// It automatically presses the ping key while aiming (holding right mouse button).
package main

import (
	"fmt"
	"os"
	"os/signal"
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

// Virtual-Key Codes
// https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
const (
	VK_RBUTTON  = 0x02 // Right mouse button
	VK_LCONTROL = 0xA2 // Left Control key
)

// keybd_event flags
const (
	KEYEVENTF_KEYUP = 0x0002
)

// Timing configuration
const (
	PingInterval   = 1 * time.Second
	PollRateActive = 50 * time.Millisecond  // Fast polling when aiming
	PollRateIdle   = 200 * time.Millisecond // Slow polling when idle
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
	fmt.Println("Press Ctrl+C to exit")

	// Graceful shutdown on Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

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
			time.Sleep(PollRateActive) // Fast polling when aiming
		} else {
			pressed = false
			time.Sleep(PollRateIdle) // Slow polling when idle
		}
	}
}
