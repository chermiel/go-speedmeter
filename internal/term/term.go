package term

func HideCursor() string {
	return "\033[?25l"
}

func ShowCursor() string {
	return "\033[?25h"
}

func ClearScreen() string {
	return "\033[H\033[2J"
}
