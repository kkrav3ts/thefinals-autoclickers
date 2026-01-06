package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kkrav3ts/thefinals-autoping/internal/keyboard"
)

func main() {
	// Graceful shutdown on Ctrl+C
	fmt.Println("Close window or press Ctrl+C to exit")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()
	keyboard.CheckLMKDelay()
}
