package main

import (
	"fmt"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkrav3ts/thefinals-autoping/internal/keyboard"
	"github.com/kkrav3ts/thefinals-autoping/internal/statistics"
)

func main() {
	fmt.Println("THE FINALS Auto-Shooting Tool. Built by Bykang.")

	// PREDEFINED INPUTS
	leftMouseButton := 0x01 // Virtual-Key Code for Left Mouse Button used as shooting key.
	mean := 72.5            // midpoint of 60-85ms delay cluster
	stdDev := 10.0          // standard deviation to create the delay cluster
	minVal := 50.0          // minimum delay
	maxVal := 105.0         // maximum delay

	// USER-BASED INPUT
	fmt.Printf("Press the key you want to use for shooting.\n")
	shotKey := keyboard.DetectKeyPress(keyboard.KeyNames)
	fmt.Printf("Auto-shooting enabled using [%s] key. Hold left mouse button to simulate repeated clicks...\n", keyboard.KeyNames[shotKey])

	// Graceful shutdown on Ctrl+C
	fmt.Println("Close window or press Ctrl+C to exit")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	// Generate a pool of delays to cycle through
	delays := statistics.GenerateClickDelays(10000, mean, stdDev, minVal, maxVal)
	delayIndex := 0

	// Infinite loop for the main process
	for {
		if keyboard.IsKeyPressed(leftMouseButton) {
			// Press key with human-like key pressed time
			keyboard.PressKey(shotKey, delays[delayIndex])
			delayIndex = (delayIndex + 1) % len(delays)

			// Human-like Delay between key pressed
			time.Sleep(delays[delayIndex])
			delayIndex = (delayIndex + 1) % len(delays)
		} else {
			// Small polling delay to avoid excessive CPU usage when idle
			time.Sleep(20 * time.Millisecond)
		}
	}
}
