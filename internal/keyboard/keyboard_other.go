//go:build !windows

package keyboard

import (
	"context"
	"errors"
	"time"
)

// ErrUnsupported is returned when an input function is called on a system
// other than Windows. The tools simulate Win32 input.
var ErrUnsupported = errors.New("keyboard input requires windows")

func IsKeyPressed(int) bool { return false }

func DetectKeyPress(context.Context, map[int]string) (int, error) {
	return 0, ErrUnsupported
}

func WaitReleased(context.Context, int) error { return ErrUnsupported }

func PressKey(int, time.Duration) error { return ErrUnsupported }

func WithTimerResolution(fn func()) { fn() }

func CheckLMKDelay(context.Context) error { return ErrUnsupported }
