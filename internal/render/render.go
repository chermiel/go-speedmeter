package render

import (
	"fmt"
	"time"

	"speedmeter/internal/format"
)

const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
	colorReset  = "\033[0m"
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

func TUI(d Data, firstDraw bool) {
	if firstDraw {
		fmt.Print("\033[H\033[2J")
		fmt.Print(colorBold)
		fmt.Println("╭──────────────────────────────────────────╮")
		fmt.Printf("│    SpeedMeter  •  %-22s│\n", d.Iface)
		fmt.Println("├──────────────────────────────────────────┤")
		fmt.Println("│  ⬇ DOWNLOAD                               │")
		fmt.Println("│  ▏                                        │")
		fmt.Println("│  ⬆ UPLOAD                                 │")
		fmt.Println("│  ▏                                        │")
		fmt.Println("├──────────────────────────────────────────┤")
		fmt.Println("│  Total  RX: --         │  TX: --          │")
		fmt.Println("│  Peak   DL: --         │  UL: --          │")
		fmt.Println("│  Runtime:  --                              │")
		fmt.Println("├──────────────────────────────────────────┤")
		fmt.Println("│  [Ctrl+C] exit                             │")
		fmt.Println("╰──────────────────────────────────────────╯")
		fmt.Print(colorReset)
		fmt.Print("\033[?25l")
	}

	pctRx := d.RxSpeed / (d.PeakRx + 1)
	pctTx := d.TxSpeed / (d.PeakTx + 1)

	// Row 4 (1-based): DL icon line → write speed + pct
	fmt.Printf("\033[4;4H%s%-14s %s%s", colorGreen, format.Speed(d.RxSpeed), format.Percent(pctRx), colorReset)
	// Row 5: DL bar
	fmt.Printf("\033[5;5H%s%s%s", colorGreen, format.BarFractional(pctRx, 38), colorReset)

	// Row 6: UL icon line → write speed + pct
	fmt.Printf("\033[6;4H%s%-14s %s%s", colorYellow, format.Speed(d.TxSpeed), format.Percent(pctTx), colorReset)
	// Row 7: UL bar
	fmt.Printf("\033[7;5H%s%s%s", colorYellow, format.BarFractional(pctTx, 38), colorReset)

	// Row 9: Total RX/TX
	fmt.Printf("\033[9;14H%-12s │  TX: %-10s", format.Bytes(d.TotalRx), format.Bytes(d.TotalTx))
	// Row 10: Peak DL/UL
	fmt.Printf("\033[10;15H%-12s │  UL: %-10s", format.Speed(d.PeakRx), format.Speed(d.PeakTx))
	// Row 11: Runtime
	fmt.Printf("\033[11;13H%-28s", d.Runtime)
}

func Plain(d Data) {
	fmt.Print("\033[H\033[2J")

	fmt.Println("==============================================")
	fmt.Printf("=  SpeedMeter — %-22s =\n", d.Iface)
	fmt.Println("==============================================")
	fmt.Printf("  DOWNLOAD: %s\n", format.Speed(d.RxSpeed))
	fmt.Printf("  UPLOAD:   %s\n", format.Speed(d.TxSpeed))
	fmt.Println("----------------------------------------------")
	fmt.Printf("  Total RX: %s\n", format.Bytes(d.TotalRx))
	fmt.Printf("  Total TX: %s\n", format.Bytes(d.TotalTx))
	fmt.Printf("  Peak DL:  %s\n", format.Speed(d.PeakRx))
	fmt.Printf("  Peak UL:  %s\n", format.Speed(d.PeakTx))
	fmt.Printf("  Runtime:   %s\n", d.Runtime)
	fmt.Println("==============================================")
}
