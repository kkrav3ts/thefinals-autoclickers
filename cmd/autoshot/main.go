package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kkrav3ts/thefinals-autoclickers/internal/app"
	"github.com/kkrav3ts/thefinals-autoclickers/internal/keyboard"
	"github.com/kkrav3ts/thefinals-autoclickers/internal/timing"
)

func main() {
	mean := flag.Float64("mean", 60, "center of the click-timing distribution, in milliseconds")
	stdDev := flag.Float64("stddev", 5, "standard deviation, in milliseconds")
	minVal := flag.Float64("min", 50, "minimum key hold and gap, in milliseconds")
	maxVal := flag.Float64("max", 80, "maximum key hold and gap, in milliseconds")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: thefinals-autoshot [flags]\n\nHold the left mouse button to repeat the shoot key you select.\nEach click holds the key for one random sample, then waits for another.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	fmt.Println("THE FINALS Auto-Shooting Tool.")
	fmt.Printf("Click timing mean %.0fms, spread %.0fms, clamped to %.0f–%.0fms.\n", *mean, *stdDev, *minVal, *maxVal)
	fmt.Println("Press the key you want to use for shooting.")

	var runErr error
	app.Run(func(ctx context.Context) {
		runErr = loop(ctx, *mean, *stdDev, *minVal, *maxVal)
	})
	if runErr != nil && runErr != context.Canceled {
		fmt.Fprintln(os.Stderr, runErr)
		os.Exit(1)
	}
}

func loop(ctx context.Context, mean, stdDev, minVal, maxVal float64) error {
	shotKey, err := keyboard.DetectKeyPress(ctx, keyboard.KeyNames)
	if err != nil {
		return err
	}
	fmt.Printf("Auto-shooting enabled using [%s]. Hold the left mouse button to fire.\n", keyboard.KeyNames[shotKey])
	fmt.Println("Close the window or press Ctrl+C to exit.")
	if err := keyboard.WaitReleased(ctx, shotKey); err != nil {
		return err
	}

	delays, err := timing.GenerateClickDelays(1000, mean, stdDev, minVal, maxVal)
	if err != nil {
		return err
	}
	fmt.Printf("Generated %d click timings.\n", len(delays))

	delayIndex := 0
	next := func() time.Duration {
		d := delays[delayIndex]
		delayIndex = (delayIndex + 1) % len(delays)
		return d
	}

	for {
		if ctx.Err() != nil {
			return nil
		}
		if !keyboard.IsKeyPressed(keyboard.VKLeftButton) {
			if !app.Sleep(ctx, 10*time.Millisecond) {
				return nil
			}
			continue
		}

		hold := next()
		gap := next()
		if err := keyboard.PressKey(shotKey, hold); err != nil {
			fmt.Fprintf(os.Stderr, "input: %v\n", err)
			if !app.Sleep(ctx, time.Second) {
				return nil
			}
			continue
		}
		if !app.Sleep(ctx, gap) {
			return nil
		}
	}
}
