package format

import (
	"fmt"
	"strings"
)

func Speed(bps float64) string {
	if bps < 0 {
		bps = 0
	}
	switch {
	case bps >= 1_000_000_000:
		return fmt.Sprintf("%7.2f Gbps", bps/1_000_000_000)
	case bps >= 1_000_000:
		return fmt.Sprintf("%7.2f Mbps", bps/1_000_000)
	case bps >= 1_000:
		return fmt.Sprintf("%7.2f Kbps", bps/1_000)
	default:
		return fmt.Sprintf("%7.2f bps", bps)
	}
}

func Bytes(b uint64) string {
	switch {
	case b >= 1_073_741_824:
		return fmt.Sprintf("%.2f GB", float64(b)/1_073_741_824)
	case b >= 1_048_576:
		return fmt.Sprintf("%.2f MB", float64(b)/1_048_576)
	case b >= 1_024:
		return fmt.Sprintf("%.2f KB", float64(b)/1_024)
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func Bar(ratio float64, width int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}
	if filled == 0 && ratio > 0 {
		filled = 1
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
