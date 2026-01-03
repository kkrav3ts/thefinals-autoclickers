// Package main provides an auto-ping utility for THE FINALS.
// It automatically presses the ping key while aiming (holding right mouse button).
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkrav3ts/thefinals-autoping/internal/keyboard"
)

// Version is set at build time via -ldflags
var Version = "dev"

// Timing configuration
const (
	PingInterval   = 1 * time.Second
	PollRateActive = 50 * time.Millisecond  // Fast polling when aiming
	PollRateIdle   = 200 * time.Millisecond // Slow polling when idle
)

func main() {
	fmt.Printf("THE FINALS Auto-Ping %s\n", Version)
	fmt.Println()

	// Prompt user to select ping key
	fmt.Print("Press the key you want to use for ping: ")
	pingKey := keyboard.DetectKeyPress()
	fmt.Printf("%s\n", keyboard.GetKeyName(pingKey))
	fmt.Println()

	// Wait for key release before starting
	keyboard.WaitForKeyRelease()

	fmt.Printf("Auto-ping enabled using [%s] key\n", keyboard.GetKeyName(pingKey))
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
		if keyboard.IsKeyPressed(keyboard.VK_RBUTTON) {
			if !pressed {
				pressed = true
				lastPing = time.Now()
				keyboard.PressKey(pingKey)
			} else if time.Since(lastPing) >= PingInterval {
				keyboard.PressKey(pingKey)
				lastPing = time.Now()
			}
			time.Sleep(PollRateActive) // Fast polling when aiming
		} else {
			pressed = false
			time.Sleep(PollRateIdle) // Slow polling when idle
		}
	}
}
