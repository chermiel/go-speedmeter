package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func isTerminal() bool {
	var ws struct{ Row, Col, Xpix, Ypix uint16 }
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	return err == 0
}

func readNetStats(iface string) (uint64, uint64, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		name := strings.TrimRight(fields[0], ":")
		if name == iface {
			rx, _ := strconv.ParseUint(fields[1], 10, 64)
			tx, _ := strconv.ParseUint(fields[9], 10, 64)
			return rx, tx, nil
		}
	}
	return 0, 0, fmt.Errorf("interface %s not found", iface)
}

func formatSpeed(bps float64) string {
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

func formatBytes(b uint64) string {
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

func bar(ratio float64, width int) string {
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

func draw(iface string, rxSpeed, txSpeed, peakRx, peakTx float64, totalRx, totalTx uint64, runtime time.Duration, firstDraw bool) {
	width := 40
	green := "\033[32m"
	yellow := "\033[33m"
	reset := "\033[0m"

	if firstDraw {
		// Only draw static layout once
		fmt.Print("\033[H\033[2J")
		fmt.Println("╔══════════════════════════════════════════╗")
		fmt.Printf("║  🌐 SpeedMeter — %-19s ║\n", iface)
		fmt.Println("╠══════════════════════════════════════════╣")
		fmt.Println("║                                          ║")
		fmt.Println("║  ⬇ DOWNLOAD:                             ║")
		fmt.Println("║                                          ║")
		fmt.Println("║                                          ║")
		fmt.Println("║  ⬆ UPLOAD:                               ║")
		fmt.Println("║                                          ║")
		fmt.Println("║                                          ║")
		fmt.Println("╠══════════════════════════════════════════╣")
		fmt.Println("║  📊 Total RX:                    ║")
		fmt.Println("║  📊 Total TX:                    ║")
		fmt.Println("║  🚀 Peak DL:                     ║")
		fmt.Println("║  🚀 Peak UL:                     ║")
		fmt.Println("║  ⏱  Runtime:                     ║")
		fmt.Println("╠══════════════════════════════════════════╣")
		fmt.Println("║  [Ctrl+C] exit                           ║")
		fmt.Println("╚══════════════════════════════════════════╝")
		fmt.Print("\033[?25l") // hide cursor
	}

	// Move to specific positions and update only the numbers
	// Row 5: Download speed
	fmt.Printf("\033[5;17H%s%-20s%s", green, formatSpeed(rxSpeed), reset)
	// Row 6: Download bar
	fmt.Printf("\033[6;5H%s%s", bar(rxSpeed/(peakRx+1), width-2), "\033[0m")
	
	// Row 8: Upload speed
	fmt.Printf("\033[8;16H%s%-20s%s", yellow, formatSpeed(txSpeed), reset)
	// Row 9: Upload bar
	fmt.Printf("\033[9;5H%s%s", bar(txSpeed/(peakTx+1), width-2), "\033[0m")
	
	// Stats
	fmt.Printf("\033[12;17H%-22s", formatBytes(totalRx))
	fmt.Printf("\033[13;17H%-22s", formatBytes(totalTx))
	fmt.Printf("\033[14;18H%-22s", formatSpeed(peakRx))
	fmt.Printf("\033[15;18H%-22s", formatSpeed(peakTx))
	fmt.Printf("\033[16;18H%-22s", runtime)
}

func drawPlain(iface string, rxSpeed, txSpeed, peakRx, peakTx float64, totalRx, totalTx uint64, runtime time.Duration) {
	fmt.Print("\033[H\033[2J")

	fmt.Println("==============================================")
	fmt.Printf("=  SpeedMeter — %-22s =\n", iface)
	fmt.Println("==============================================")
	fmt.Printf("  DOWNLOAD: %s\n", formatSpeed(rxSpeed))
	fmt.Printf("  UPLOAD:   %s\n", formatSpeed(txSpeed))
	fmt.Println("----------------------------------------------")
	fmt.Printf("  Total RX: %s\n", formatBytes(totalRx))
	fmt.Printf("  Total TX: %s\n", formatBytes(totalTx))
	fmt.Printf("  Peak DL:  %s\n", formatSpeed(peakRx))
	fmt.Printf("  Peak UL:  %s\n", formatSpeed(peakTx))
	fmt.Printf("  Runtime:   %s\n", runtime)
	fmt.Println("==============================================")
}

func main() {
	iface := "eth0"
	if len(os.Args) > 1 {
		iface = os.Args[1]
	}

	useColor := isTerminal()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var rxSpeed, txSpeed, peakRx, peakTx float64
	var totalRx, totalTx uint64

	lastRx, lastTx, err := readNetStats(iface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if useColor {
		fmt.Print("\033[?25l")
		defer fmt.Print("\033[?25h")
	}

	startTime := time.Now()
	var tickCount int
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			return
		case <-ticker.C:
			rx, tx, err := readNetStats(iface)
			if err != nil {
				continue
			}

			deltaRx := rx - lastRx
			deltaTx := tx - lastTx

			rxSpeed = float64(deltaRx) * 8
			txSpeed = float64(deltaTx) * 8
			totalRx = rx
			totalTx = tx

			if rxSpeed > peakRx {
				peakRx = rxSpeed
			}
			if txSpeed > peakTx {
				peakTx = txSpeed
			}

			lastRx = rx
			lastTx = tx

			runtime := time.Since(startTime).Round(time.Second)

			if useColor {
				draw(iface, rxSpeed, txSpeed, peakRx, peakTx, totalRx, totalTx, runtime, tickCount == 0)
			} else {
				drawPlain(iface, rxSpeed, txSpeed, peakRx, peakTx, totalRx, totalTx, runtime)
			}
			tickCount++
		}
	}
}
