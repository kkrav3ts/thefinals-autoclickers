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

// WaitForKeyRelease waits until all keys are released.
func WaitForKeyRelease() {
	for {
		anyPressed := false
		for vk := 0x08; vk <= 0xFE; vk++ {
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

		// Now wait for a key press (skip mouse buttons 0x01-0x07)
		var pressedKey int
		for pressedKey == 0 {
			for vk := 0x08; vk <= 0xFE; vk++ {
				if IsKeyPressed(vk) {
					pressedKey = vk
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}

		// Check if the key is supported
		if _, ok := keyNames[pressedKey]; ok {
			return pressedKey
		}

		fmt.Printf("Key 0x%02X is not supported. Please press another key: ", pressedKey)
	}
}
