package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kkrav3ts/thefinals-autoclickers/internal/app"
	"github.com/kkrav3ts/thefinals-autoclickers/internal/keyboard"
)

func main() {
	fmt.Println("THE FINALS Delay Checker.")
	fmt.Println("Close the window or press Ctrl+C to exit.")

	var runErr error
	app.Run(func(ctx context.Context) {
		runErr = keyboard.CheckLMKDelay(ctx)
	})
	if runErr != nil && runErr != context.Canceled {
		fmt.Fprintln(os.Stderr, runErr)
		os.Exit(1)
	}
}
