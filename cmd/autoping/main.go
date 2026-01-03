package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkrav3ts/thefinals-autoping/internal/keyboard"
)

func main() {
	fmt.Println("THE FINALS Auto-Ping Tool.")

	// PREDEFINED INPUTS
	aimKey := 0x02 // Virtual-Key Code for Right Mouse Button used as aiming key.
	PingInterval := 1 * time.Second
	PollRateActive := 100 * time.Millisecond // Polling when aiming
	PollRateIdle := 200 * time.Millisecond   // Slow polling when idle

	// USER-BASED INPUT
	fmt.Printf("Press the key you want to use for ping.\n")
	pingKey := keyboard.DetectKeyPress(keyboard.KeyNames)
	fmt.Printf("Auto-ping enabled using [%s] key. Start aiming with right mouse button...\n", keyboard.KeyNames[pingKey])

	// Graceful shutdown on Ctrl+C
	fmt.Println("Close window or press Ctrl+C to exit")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	// Infinite loop for the process with variable polling rate
	var nextPingTime time.Time
	for {
		if keyboard.IsKeyPressed(aimKey) {
			now := time.Now()
			if nextPingTime.IsZero() || now.After(nextPingTime) {
				keyboard.PressKey(pingKey, 10*time.Millisecond)
				nextPingTime = now.Add(PingInterval) // Calculate next ping time
			}
			time.Sleep(PollRateActive)
		} else {
			nextPingTime = time.Time{} // Reset when not pressing
			time.Sleep(PollRateIdle)
		}
	}
}
