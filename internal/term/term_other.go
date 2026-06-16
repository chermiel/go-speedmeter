//go:build !linux

package term

func IsTerminal() bool {
	return false
}
