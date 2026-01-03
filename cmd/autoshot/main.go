package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// IsKeyPressed returns true if the specified virtual key is currently pressed.
func IsKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

// CheckLMKDelay prints delay between Left Mouse Press and Release
func CheckLMKDelay() {
	// Left mouse button virtual key code
	leftMouseButton := 0x01 // VK_LBUTTON

	var isPressed bool
	var pressTime time.Time

	fmt.Println("Monitoring left mouse button clicks...")
	fmt.Println("Press and release the left mouse button to see the delay.")

	// Poll for mouse button state changes
	for {
		currentState := IsKeyPressed(leftMouseButton)

		// Detect button press (transition from not pressed to pressed)
		if currentState && !isPressed {
			pressTime = time.Now()
			isPressed = true
		}

		// Detect button release (transition from pressed to not pressed)
		if !currentState && isPressed {
			releaseTime := time.Now()
			delay := releaseTime.Sub(pressTime)
			fmt.Printf("Click delay: %d ms\n", delay.Milliseconds())
			isPressed = false
		}

		// Small delay to avoid excessive CPU usage
		time.Sleep(10 * time.Millisecond)
	}
}

// GenerateClickDelays generates a list of realistic click delays.
// The delays follow a normal distribution clustered around 60-85ms, with bounds of 50-105ms.
//
// Parameters:
//   - count: number of delays to generate
//
// Returns: slice of time.Duration representing delays (already multiplied by time.Millisecond)
func GenerateClickDelays(count int) []time.Duration {
	if count <= 0 {
		return []time.Duration{}
	}

	delays := make([]time.Duration, count)

	// Distribution parameters
	const (
		mean   = 72.5  // midpoint of 60-85ms cluster
		stdDev = 10.0  // standard deviation to create the cluster
		minVal = 50.0  // minimum delay
		maxVal = 105.0 // maximum delay
	)

	for i := 0; i < count; i++ {
		// Generate value from normal distribution
		val := mean + stdDev*rand.NormFloat64()

		// Clamp to bounds [50, 105]
		val = math.Max(minVal, math.Min(maxVal, val))

		// Convert to time.Duration by multiplying milliseconds
		delays[i] = time.Duration(math.Round(val)) * time.Millisecond
	}

	return delays
}

// PressKey simulates a key press and release for the specified virtual key.
func PressKey(vk int) {
	// Press key down
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, 0, 0)
	// Release key
	keyEventFKeyUp := 0x0002
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, uintptr(keyEventFKeyUp), 0)
}

func main() {
	leftMouseButton := 0x01 // VK_LBUTTON
	keyL := 0x4C            // Virtual key code for 'L'

	var wasPressed bool

	fmt.Println("Hold left mouse button to simulate repeated L key presses...")
	fmt.Println("Press Ctrl+C to exit")

	// Graceful shutdown on Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	// Generate a pool of delays to cycle through
	delays := GenerateClickDelays(1000)
	delayIndex := 0

	for {
		currentState := IsKeyPressed(leftMouseButton)

		// Detect button press (transition from not pressed to pressed)
		if currentState && !wasPressed {
			fmt.Println("Left mouse button pressed - starting simulation...")
			wasPressed = true
			delayIndex = 0 // Reset delay index for new press cycle

			// Step 1: Wait initial delay before first key press
			if len(delays) > 0 {
				time.Sleep(delays[delayIndex])
				delayIndex = (delayIndex + 1) % len(delays)
			}

			// Loop steps 2-3: Press L -> Wait delay -> Press L -> Wait delay... until released
			for {
				// Check if physical button is still pressed
				time.Sleep(5 * time.Millisecond) // Small delay to let any previous events settle
				if !IsKeyPressed(leftMouseButton) {
					break // Physical button was released, exit loop
				}

				// Step 2: Press virtual keyboard key L
				//PressKey(keyL)
				// Press key down
				_, _, _ = keybdEvent.Call(uintptr(keyL), 0, 0, 0)

				// Step 3: Wait delay before next key press
				if len(delays) > 0 {
					time.Sleep(delays[delayIndex])
					delayIndex = (delayIndex + 1) % len(delays)
				}

				// Release key
				keyEventFKeyUp := 0x0002
				_, _, _ = keybdEvent.Call(uintptr(keyL), 0, uintptr(keyEventFKeyUp), 0)

				// Step 3: Wait delay before next key press
				if len(delays) > 0 {
					time.Sleep(delays[delayIndex])
					delayIndex = (delayIndex + 1) % len(delays)
				}
			}

			// Button was released during the loop
			fmt.Println("Left mouse button released - stopping simulation")
			wasPressed = false
		}

		// Small polling delay to avoid excessive CPU usage when idle
		time.Sleep(10 * time.Millisecond)
	}
}
