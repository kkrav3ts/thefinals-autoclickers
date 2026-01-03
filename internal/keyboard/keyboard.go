//go:build windows

// Package keyboard provides Windows keyboard input handling via user32.dll.
package keyboard

import (
	"fmt"
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybdEvent       = user32.NewProc("keybd_event")
)

// keybd_event flags
const (
	keyEventFKeyUp = 0x0002
)

// IsKeyPressed returns true if the specified virtual key is currently pressed.
func IsKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

// PressKey simulates a key press and release for the specified virtual key.
func PressKey(vk int) {
	keybdEvent.Call(uintptr(vk), 0, 0, 0)
	keybdEvent.Call(uintptr(vk), 0, uintptr(keyEventFKeyUp), 0)
}

// WaitForKeyRelease waits until all supported keys are released.
func WaitForKeyRelease() {
	for {
		anyPressed := false
		for _, vk := range supportedKeys {
			if IsKeyPressed(vk) {
				anyPressed = true
				break
			}
		}
		if !anyPressed {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// DetectKeyPress waits for the user to press a supported key and returns its virtual key code.
func DetectKeyPress() int {
	for {
		// Wait for all keys to be released first
		WaitForKeyRelease()

		// Wait for a supported key press
		for {
			for _, vk := range supportedKeys {
				if IsKeyPressed(vk) {
					return vk
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// DetectAnyKeyPress waits for any key press and returns true if it's supported.
// Returns the virtual key code and whether it's a supported key.
func DetectAnyKeyPress() (vk int, supported bool) {
	WaitForKeyRelease()

	// Scan all possible keys (0x08-0xFE, skipping mouse buttons)
	for {
		for k := 0x08; k <= 0xFE; k++ {
			if IsKeyPressed(k) {
				return k, IsSupportedKey(k)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// PromptForKey prompts the user and waits for a supported key press.
// Keeps prompting if an unsupported key is pressed.
func PromptForKey() int {
	for {
		vk, supported := DetectAnyKeyPress()
		if supported {
			return vk
		}
		fmt.Printf("Key 0x%02X is not supported. Please press another key: ", vk)
	}
}
