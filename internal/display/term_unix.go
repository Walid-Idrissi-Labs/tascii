//go:build !windows

package display

import (
	"syscall"
	"unsafe"
)

// winsize mirrors the kernel's struct winsize used by TIOCGWINSZ.
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// getTermWidth calls the TIOCGWINSZ ioctl to read the terminal's column count.
// Returns 0 if the call fails or the result is unusable (e.g. piped output).
func getTermWidth() int {
	ws := &winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 || ws.Col < 20 {
		return 0
	}
	return int(ws.Col)
}