package keyboard

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybdEvent       = user32.NewProc("keybd_event")
)

// KeyNames maps virtual key codes to human-readable names.
var KeyNames = map[int]string{
	0x08: "Backspace", 0x09: "Tab", 0x0D: "Enter",
	0x13: "Pause", 0x14: "Caps Lock", 0x1B: "Escape", 0x20: "Space",
	0x21: "Page Up", 0x22: "Page Down", 0x23: "End", 0x24: "Home",
	0x25: "Left Arrow", 0x26: "Up Arrow", 0x27: "Right Arrow", 0x28: "Down Arrow",
	0x2D: "Insert", 0x2E: "Delete",
	0x30: "0", 0x31: "1", 0x32: "2", 0x33: "3", 0x34: "4", 0x35: "5", 0x36: "6", 0x37: "7", 0x38: "8", 0x39: "9",

	0x41: "A", 0x42: "B", 0x43: "C", 0x44: "D", 0x45: "E", 0x46: "F", 0x47: "G", 0x48: "H", 0x49: "I", 0x4A: "J",
	0x4B: "K", 0x4C: "L", 0x4D: "M", 0x4E: "N", 0x4F: "O", 0x50: "P", 0x51: "Q", 0x52: "R", 0x53: "S", 0x54: "T",
	0x55: "U", 0x56: "V", 0x57: "W", 0x58: "X", 0x59: "Y", 0x5A: "Z",

	0x60: "Numpad 0", 0x61: "Numpad 1", 0x62: "Numpad 2", 0x63: "Numpad 3", 0x64: "Numpad 4", 0x65: "Numpad 5",
	0x66: "Numpad 6", 0x67: "Numpad 7", 0x68: "Numpad 8", 0x69: "Numpad 9",
	0x6A: "Numpad *", 0x6B: "Numpad +", 0x6D: "Numpad -", 0x6E: "Numpad .", 0x6F: "Numpad /",

	0x70: "F1", 0x71: "F2", 0x72: "F3", 0x73: "F4", 0x74: "F5", 0x75: "F6", 0x76: "F7", 0x77: "F8", 0x78: "F9",
	0x79: "F10", 0x7A: "F11", 0x7B: "F12",

	0xA0: "Left Shift", 0xA1: "Right Shift", 0xA2: "Left Ctrl", 0xA3: "Right Ctrl", 0xA4: "Left Alt", 0xA5: "Right Alt",

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

// PressKey simulates a key press and release for the specified virtual key with definable delay.
func PressKey(vk int, delay time.Duration) {
	// Press Key Down
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, 0, 0)

	// Keep Key Pressed
	time.Sleep(delay)

	// Release Key
	keyEventFKeyUp := 0x0002
	_, _, _ = keybdEvent.Call(uintptr(vk), 0, uintptr(keyEventFKeyUp), 0)
}

// ====== DELAY CHECKER ======
var (
	winmm           = syscall.NewLazyDLL("winmm.dll")
	timeBeginPeriod = winmm.NewProc("timeBeginPeriod")
	timeEndPeriod   = winmm.NewProc("timeEndPeriod")
)

// setHighResTimer sets Windows timer resolution to 1ms for accurate Sleep calls.
// Returns a cleanup function that should be deferred.
func setHighResTimer() func() {
	timeBeginPeriod.Call(1)
	return func() {
		timeEndPeriod.Call(1)
	}
}

// CheckLMKDelay prints delay between Left Mouse Press and Release, and between consecutive presses.
// Uses high-resolution timing and optimized polling for accurate measurements.
func CheckLMKDelay() {
	// Lock this goroutine to a single OS thread to prevent scheduling jitter
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Set Windows timer resolution to 1ms (default is ~15.6ms)
	cleanup := setHighResTimer()
	defer cleanup()

	// Left mouse button virtual key code
	const leftMouseButton = 0x01 // VK_LBUTTON

	var isPressed bool
	var pressTime time.Time
	var lastPressTime time.Time
	var clickCount int

	fmt.Println("Monitoring left mouse button clicks (high-precision mode)...")
	fmt.Println("Press and release the left mouse button to see delays.")
	// Pre-fetch the key state check to reduce syscall overhead in hot path	fmt.Println()
	checkKey := func() bool {
		ret, _, _ := getAsyncKeyState.Call(uintptr(leftMouseButton))
		return ret&0x8000 != 0
	}

	// Poll for mouse button state changes
	for {
		currentState := checkKey()

		// Detect button press (transition from not pressed to pressed)
		if currentState && !isPressed {
			pressTime = time.Now()
			isPressed = true
			clickCount++

			// Calculate interval between presses (if not the first press)
			if !lastPressTime.IsZero() {
				interval := pressTime.Sub(lastPressTime)
				fmt.Printf("[Click #%d] Interval since last press: %d ms (%.2f ms)\n",
					clickCount, interval.Milliseconds(), float64(interval.Microseconds())/1000.0)
			} else {
				fmt.Printf("[Click #%d] First press detected\n", clickCount)
			}
		}

		// Detect button release (transition from pressed to not pressed)
		if !currentState && isPressed {
			releaseTime := time.Now()
			holdDuration := releaseTime.Sub(pressTime)
			fmt.Printf("[Click #%d] Hold duration: %d ms (%.2f ms)\n",
				clickCount, holdDuration.Milliseconds(), float64(holdDuration.Microseconds())/1000.0)
			fmt.Println()

			lastPressTime = pressTime // Save for next interval calculation
			isPressed = false
		}

		// Minimal sleep to balance CPU usage and accuracy
		// With timeBeginPeriod(1), this actually sleeps ~1ms instead of ~15ms
		time.Sleep(500 * time.Microsecond)
	}
}
