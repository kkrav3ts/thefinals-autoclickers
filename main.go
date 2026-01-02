package main

import (
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	keybd_event      = user32.NewProc("keybd_event")
)

const (
	VK_RBUTTON      = 0x02
	VK_LCONTROL     = 0xA2
	KEYEVENTF_KEYUP = 0x0002
)

func isKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

func pressKey(vk int) {
	keybd_event.Call(uintptr(vk), 0, 0, 0)
	keybd_event.Call(uintptr(vk), 0, uintptr(KEYEVENTF_KEYUP), 0)
}

func main() {
	var pressed bool
	var t time.Time

	for {
		state := isKeyPressed(VK_RBUTTON)

		if state {
			if !pressed {
				pressed = true
				t = time.Now()
				pressKey(VK_LCONTROL)
			} else if time.Since(t) >= 1*time.Second {
				pressKey(VK_LCONTROL)
				t = time.Now()
			}
		} else {
			pressed = false
		}

		time.Sleep(100 * time.Millisecond)
	}
}
