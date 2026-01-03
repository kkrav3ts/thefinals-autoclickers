package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkrav3ts/thefinals-autoping/internal/keyboard"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybdEvent       = user32.NewProc("keybd_event")
)

func main() {
	fmt.Println("THE FINALS Auto-Ping Tool. Built by Bykang.")

	// Prompt user to select ping key.
	fmt.Printf("Press the key you want to use for ping.\n")
	pingKey := DetectKeyPress(keyNames)
	fmt.Printf("Auto-ping enabled using [%s] key. Good luck, contestant!\n", keyNames[pingKey])

	// Define aiming Key Virtual-Key Code
	aimKey := 0x02 // Right mouse button

	// Timing configuration
	PingInterval := 1 * time.Second
	PollRateActive := 100 * time.Millisecond // Polling when aiming
	PollRateIdle := 200 * time.Millisecond   // Slow polling when idle

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
		if IsKeyPressed(aimKey) {
			now := time.Now()
			if nextPingTime.IsZero() || now.After(nextPingTime) {
				PressKey(pingKey)
				nextPingTime = now.Add(PingInterval) // Calculate next ping time
			}
			time.Sleep(PollRateActive)
		} else {
			nextPingTime = time.Time{} // Reset when not pressing
			time.Sleep(PollRateIdle)
		}
	}
}
