//go:build windows

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
	pingKey := keyboard.PromptForKey()
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

	// Use ticker for accurate ping interval timing
	pingTicker := time.NewTicker(PingInterval)
	defer pingTicker.Stop()

	var aiming bool

	for {
		if keyboard.IsKeyPressed(keyboard.VK_RBUTTON) {
			if !aiming {
				// Just started aiming - ping immediately
				aiming = true
				keyboard.PressKey(pingKey)
				pingTicker.Reset(PingInterval)
			}
			// Check ticker for subsequent pings
			select {
			case <-pingTicker.C:
				keyboard.PressKey(pingKey)
			default:
			}
			time.Sleep(PollRateActive)
		} else {
			aiming = false
			time.Sleep(PollRateIdle)
		}
	}
}
