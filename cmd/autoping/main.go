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
	interval := flag.Duration("interval", time.Second, "time between pings while aiming")
	hold := flag.Duration("hold", 10*time.Millisecond, "how long the ping key is held")
	mean := flag.Float64("mean", 80, "extra jitter around each ping, in milliseconds")
	stdDev := flag.Float64("stddev", 5, "jitter standard deviation, in milliseconds")
	minVal := flag.Float64("min", 70, "minimum jitter, in milliseconds")
	maxVal := flag.Float64("max", 100, "maximum jitter, in milliseconds")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: thefinals-autoping [flags]\n\nHold right mouse button to tap the ping key you select.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *interval < 0 || *hold < 0 {
		fmt.Fprintln(os.Stderr, "interval and hold must be zero or positive")
		os.Exit(2)
	}

	fmt.Println("THE FINALS Auto-Ping Tool.")
	fmt.Printf("Ping spacing %s, hold %s, jitter mean %.0fms (%.0f–%.0fms).\n", *interval, *hold, *mean, *minVal, *maxVal)
	fmt.Println("Press the key you want to use for ping.")

	var runErr error
	app.Run(func(ctx context.Context) {
		runErr = loop(ctx, *interval, *hold, *mean, *stdDev, *minVal, *maxVal)
	})
	if runErr != nil && runErr != context.Canceled {
		fmt.Fprintln(os.Stderr, runErr)
		os.Exit(1)
	}
}

func loop(ctx context.Context, interval, hold time.Duration, mean, stdDev, minVal, maxVal float64) error {
	pingKey, err := keyboard.DetectKeyPress(ctx, keyboard.KeyNames)
	if err != nil {
		return err
	}
	fmt.Printf("Auto-ping enabled using [%s]. Aim with the right mouse button.\n", keyboard.KeyNames[pingKey])
	fmt.Println("Close the window or press Ctrl+C to exit.")
	if err := keyboard.WaitReleased(ctx, pingKey); err != nil {
		return err
	}

	delays, err := timing.GenerateClickDelays(500, mean, stdDev, minVal, maxVal)
	if err != nil {
		return err
	}
	fmt.Printf("Generated %d jitter values.\n", len(delays))

	delayIndex := 0
	next := func() time.Duration {
		d := delays[delayIndex]
		delayIndex = (delayIndex + 1) % len(delays)
		return d
	}

	var nextPing time.Time
	for {
		if ctx.Err() != nil {
			return nil
		}
		if keyboard.IsKeyPressed(keyboard.VKRightButton) {
			now := time.Now()
			if nextPing.IsZero() {
				nextPing = now.Add(next())
			}
			if !now.Before(nextPing) {
				if err := keyboard.PressKey(pingKey, hold); err != nil {
					fmt.Fprintf(os.Stderr, "input: %v\n", err)
					if !app.Sleep(ctx, time.Second) {
						return nil
					}
					continue
				}
				nextPing = time.Now().Add(interval + next())
			}
			if !app.Sleep(ctx, 10*time.Millisecond) {
				return nil
			}
			continue
		}
		nextPing = time.Time{}
		if !app.Sleep(ctx, 200*time.Millisecond) {
			return nil
		}
	}
}
