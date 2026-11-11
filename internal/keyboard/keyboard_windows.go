//go:build windows

package keyboard

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	inputKeyboard     = 1
	keyeventfKeyup    = 0x0002
	keyeventfScancode = 0x0008
	keyeventfExtended = 0x0001
	mapvkVkToVsc      = 0
)

// keybdInput matches KEYBDINPUT. extraInfo is pointer-sized, so the struct
// is 24 bytes on 64-bit Windows.
type keybdInput struct {
	vk        uint16
	scan      uint16
	flags     uint32
	time      uint32
	extraInfo uintptr
}

// kinput matches INPUT for a keyboard event on 64-bit Windows.
// The type field is followed by padding, then the 32-byte input union.
// KEYBDINPUT only fills 24 of those bytes, so pad keeps sizeof(INPUT) at 40.
type kinput struct {
	typ uint32
	ki  keybdInput
	pad uint64
}

// Fail compilation unless this matches sizeof(INPUT) on 64-bit Windows.
var _ = [1]byte{}[unsafe.Sizeof(kinput{})-40]

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	winmm            = windows.NewLazySystemDLL("winmm.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	sendInput        = user32.NewProc("SendInput")
	mapVirtualKeyW   = user32.NewProc("MapVirtualKeyW")
	timeBeginPeriod  = winmm.NewProc("timeBeginPeriod")
	timeEndPeriod    = winmm.NewProc("timeEndPeriod")
)

// IsKeyPressed reports whether vk is currently down.
func IsKeyPressed(vk int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

// DetectKeyPress waits until every watched key is up, then returns the next
// one that goes down. Ignoring keys that are already held avoids binding
// whatever was used to focus the console.
func DetectKeyPress(ctx context.Context, keys map[int]string) (int, error) {
	watched := ordered(keys)
	for anyDown(watched) {
		if !wait(ctx, 20*time.Millisecond) {
			return 0, ctx.Err()
		}
	}
	for {
		for _, vk := range watched {
			if IsKeyPressed(vk) {
				return vk, nil
			}
		}
		if !wait(ctx, 20*time.Millisecond) {
			return 0, ctx.Err()
		}
	}
}

// WaitReleased blocks until vk is up.
func WaitReleased(ctx context.Context, vk int) error {
	for IsKeyPressed(vk) {
		if !wait(ctx, 20*time.Millisecond) {
			return ctx.Err()
		}
	}
	return nil
}

// PressKey synthesizes a key down and key up using SendInput scan codes.
// keybd_event is ignored by many current games; scan-code input is not.
func PressKey(vk int, hold time.Duration) error {
	if err := keyEvent(vk, false); err != nil {
		return fmt.Errorf("key down: %w", err)
	}
	time.Sleep(hold)
	if err := keyEvent(vk, true); err != nil {
		return fmt.Errorf("key up: %w", err)
	}
	return nil
}

// WithTimerResolution requests a 1ms sleep resolution for fn.
// The default Windows resolution is about 15.6ms, which stretches the
// short delays these tools generate.
func WithTimerResolution(fn func()) {
	timeBeginPeriod.Call(1)
	defer timeEndPeriod.Call(1)
	fn()
}

// CheckLMKDelay prints hold duration and the gap between left-button presses.
func CheckLMKDelay(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Println("Monitoring left mouse button clicks (high-precision mode)...")
	fmt.Println("Press and release the left mouse button to see delays.")
	fmt.Println()

	var (
		isPressed     bool
		pressTime     time.Time
		lastPressTime time.Time
		clickCount    int
		totalHold     time.Duration
		totalInterval time.Duration
		intervals     int
	)

	for {
		if ctx.Err() != nil {
			return nil
		}
		down := IsKeyPressed(VKLeftButton)

		if down && !isPressed {
			pressTime = time.Now()
			isPressed = true
			clickCount++
			if !lastPressTime.IsZero() {
				interval := pressTime.Sub(lastPressTime)
				intervals++
				totalInterval += interval
				fmt.Printf("[Click #%d] Interval since last press: %.2f ms (avg %.2f ms)\n",
					clickCount, ms(interval), ms(totalInterval)/float64(intervals))
			} else {
				fmt.Printf("[Click #%d] First press detected\n", clickCount)
			}
		}

		if !down && isPressed {
			hold := time.Since(pressTime)
			totalHold += hold
			fmt.Printf("[Click #%d] Hold duration: %.2f ms (avg %.2f ms)\n\n",
				clickCount, ms(hold), ms(totalHold)/float64(clickCount))
			lastPressTime = pressTime
			isPressed = false
		}

		if !wait(ctx, 500*time.Microsecond) {
			return nil
		}
	}
}

func keyEvent(vk int, up bool) error {
	scan, _, _ := mapVirtualKeyW.Call(uintptr(vk), mapvkVkToVsc)

	var in kinput
	in.typ = inputKeyboard
	var flags uint32
	if up {
		flags |= keyeventfKeyup
	}
	if scan != 0 {
		in.ki.scan = uint16(scan)
		flags |= keyeventfScancode
		if extendedScan(vk) {
			flags |= keyeventfExtended
		}
	} else {
		in.ki.vk = uint16(vk)
	}
	in.ki.flags = flags

	sent, _, err := sendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
	if sent != 1 {
		if err != nil {
			return err
		}
		return fmt.Errorf("SendInput accepted 0 events")
	}
	return nil
}

func ordered(keys map[int]string) []int {
	out := make([]int, 0, len(keys))
	for vk := range keys {
		out = append(out, vk)
	}
	slices.Sort(out)
	return out
}

func anyDown(keys []int) bool {
	for _, vk := range keys {
		if IsKeyPressed(vk) {
			return true
		}
	}
	return false
}

func wait(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func ms(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}
