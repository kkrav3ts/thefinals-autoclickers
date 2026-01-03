package main

import (
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybdEvent       = user32.NewProc("keybd_event")
)

// keyNames maps virtual key codes to human-readable names.
var keyNames = map[int]string{
	0x08: "Backspace", 0x09: "Tab", 0x0D: "Enter",
	0x13: "Pause", 0x14: "Caps Lock", 0x1B: "Escape", 0x20: "Space",
	0x21: "Page Up", 0x22: "Page Down", 0x23: "End", 0x24: "Home",
	0x25: "Left Arrow", 0x26: "Up Arrow", 0x27: "Right Arrow", 0x28: "Down Arrow",
	0x2D: "Insert", 0x2E: "Delete",
	0x30: "0", 0x31: "1", 0x32: "2", 0x33: "3", 0x34: "4",
	0x35: "5", 0x36: "6", 0x37: "7", 0x38: "8", 0x39: "9",
	0x41: "A", 0x42: "B", 0x43: "C", 0x44: "D", 0x45: "E", 0x46: "F", 0x47: "G",
	0x48: "H", 0x49: "I", 0x4A: "J", 0x4B: "K", 0x4C: "L", 0x4D: "M", 0x4E: "N",
	0x4F: "O", 0x50: "P", 0x51: "Q", 0x52: "R", 0x53: "S", 0x54: "T", 0x55: "U",
	0x56: "V", 0x57: "W", 0x58: "X", 0x59: "Y", 0x5A: "Z",
	0x60: "Numpad 0", 0x61: "Numpad 1", 0x62: "Numpad 2", 0x63: "Numpad 3",
	0x64: "Numpad 4", 0x65: "Numpad 5", 0x66: "Numpad 6", 0x67: "Numpad 7",
	0x68: "Numpad 8", 0x69: "Numpad 9",
	0x6A: "Numpad *", 0x6B: "Numpad +", 0x6D: "Numpad -", 0x6E: "Numpad .", 0x6F: "Numpad /",
	0x70: "F1", 0x71: "F2", 0x72: "F3", 0x73: "F4", 0x74: "F5", 0x75: "F6",
	0x76: "F7", 0x77: "F8", 0x78: "F9", 0x79: "F10", 0x7A: "F11", 0x7B: "F12",
	0xA0: "Left Shift", 0xA1: "Right Shift", 0xA2: "Left Ctrl", 0xA3: "Right Ctrl",
	0xA4: "Left Alt", 0xA5: "Right Alt",
	0xBA: ";", 0xBB: "=", 0xBC: ",", 0xBD: "-", 0xBE: ".", 0xBF: "/", 0xC0: "`",
	0xDB: "[", 0xDC: "\\", 0xDD: "]", 0xDE: "'",
}

// IsKeyPressed returns true if the specified virtual key is currently pressed.
func IsKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

// DetectKeyPress waits for any key press and returns the virtual key code.
func DetectKeyPress(keysList map[int]string) int {
	// Infinite loop for the process with 50ms pause
	for {
		// Poll each provided key for pressed state
		for k := range keysList {
			if IsKeyPressed(k) {
				return k
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// PressKey simulates a key press and release for the specified virtual key.
func PressKey(vk int) {
	// Press key down
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, 0, 0)
	// Release key
	keyEventFKeyUp := 0x0002
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, uintptr(keyEventFKeyUp), 0)
}
