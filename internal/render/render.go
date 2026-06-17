package render

import (
	"fmt"
	"strings"
	"time"

	"speedmeter/internal/format"
)

const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorReset  = "\033[0m"
	barWidth    = 40
	innerWidth  = 50 // lebar isi, di luar border kiri-kanan
)

type Data struct {
	Iface   string
	RxSpeed float64
	TxSpeed float64
	PeakRx  float64
	PeakTx  float64
	TotalRx uint64
	TotalTx uint64
	Runtime time.Duration
}

// ratio aman dari div-by-zero, hasil 0..1
func ratio(v, peak float64) float64 {
	if peak <= 0 {
		return 0
	}
	if v > peak {
		return 1
	}
	return v / peak
}

// bar bikin progress bar mulus pakai partial block
func bar(pct float64, w int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	blocks := []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}
	total := pct * float64(w)
	full := int(total)

	var b strings.Builder
	b.WriteString(strings.Repeat("█", full))
	if full < w {
		frac := total - float64(full)
		b.WriteRune(blocks[int(frac*8)])
		b.WriteString(strings.Repeat("░", w-full-1))
	}
	return b.String()
}

// padR rune-aware, aman buat karakter Unicode
func padR(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return string(r[:n])
	}
	return s + strings.Repeat(" ", n-len(r))
}

func row(s string) string { return "│ " + padR(s, innerWidth) + " │" }
func divTop() string      { return "┌" + strings.Repeat("─", innerWidth+2) + "┐" }
func divMid() string      { return "├" + strings.Repeat("─", innerWidth+2) + "┤" }
func divBot() string      { return "└" + strings.Repeat("─", innerWidth+2) + "┘" }

func TUI(d Data, firstDraw bool) {
	pctRx := ratio(d.RxSpeed, d.PeakRx)
	pctTx := ratio(d.TxSpeed, d.PeakTx)

	lines := []string{
		colorBold + colorCyan + divTop(),
		row("SpeedMeter   " + d.Iface),
		divMid() + colorReset,
		colorGreen + row(fmt.Sprintf("DL  %-14s  %s", format.Speed(d.RxSpeed), format.Percent(pctRx))),
		colorGreen + row("    "+bar(pctRx, barWidth)),
		colorYellow + row(fmt.Sprintf("UL  %-14s  %s", format.Speed(d.TxSpeed), format.Percent(pctTx))),
		colorYellow + row("    "+bar(pctTx, barWidth)),
		colorCyan + divMid() + colorReset,
		row(fmt.Sprintf("RX %-18s TX %s", format.Bytes(d.TotalRx), format.Bytes(d.TotalTx))),
		row(fmt.Sprintf("Peak DL %-11s Peak UL %s", format.Speed(d.PeakRx), format.Speed(d.PeakTx))),
		row(fmt.Sprintf("Uptime %s", d.Runtime.Round(time.Second))),
		colorCyan + divBot(),
		colorDim + "  ctrl+c to exit" + colorReset,
	}

	if firstDraw {
		fmt.Print("\033[H\033[2J\033[?25l")
	} else {
		fmt.Printf("\033[%dA", len(lines)) // naik sesuai jumlah baris, no more magic number
	}
	fmt.Println(strings.Join(lines, "\n"))
}

func Plain(d Data) {
	fmt.Print("\033[H\033[2J")
	fmt.Println(divTop())
	fmt.Println(row("SpeedMeter   " + d.Iface))
	fmt.Println(divMid())
	fmt.Printf("  DOWNLOAD : %s\n", format.Speed(d.RxSpeed))
	fmt.Printf("  UPLOAD   : %s\n", format.Speed(d.TxSpeed))
	fmt.Println(divMid())
	fmt.Printf("  Total RX : %s\n", format.Bytes(d.TotalRx))
	fmt.Printf("  Total TX : %s\n", format.Bytes(d.TotalTx))
	fmt.Printf("  Peak DL  : %s\n", format.Speed(d.PeakRx))
	fmt.Printf("  Peak UL  : %s\n", format.Speed(d.PeakTx))
	fmt.Printf("  Runtime  : %s\n", d.Runtime.Round(time.Second))
	fmt.Println(divBot())
}
