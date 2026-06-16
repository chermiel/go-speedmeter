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

// BarFractional renders a progress bar using 1/8th block characters for smoother visuals.
func BarFractional(ratio float64, width int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	blocks := []string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}

	totalEighths := ratio * float64(width) * 8
	fullBlocks := int(totalEighths) / 8
	remainder := int(totalEighths) % 8

	var result strings.Builder
	for i := 0; i < width; i++ {
		switch {
		case i < fullBlocks:
			result.WriteString("█")
		case i == fullBlocks && remainder > 0:
			result.WriteString(blocks[remainder])
		default:
			result.WriteString(" ")
		}
	}
	return result.String()
}

// Percent formats a ratio as a right-aligned percentage string.
func Percent(ratio float64) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return fmt.Sprintf("%5.1f%%", ratio*100)
}
