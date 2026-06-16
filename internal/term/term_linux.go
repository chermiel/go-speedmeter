//go:build linux

package term

import (
	"syscall"
	"unsafe"
)

func IsTerminal() bool {
	var ws struct{ Row, Col, Xpix, Ypix uint16 }
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	return err == 0
}
