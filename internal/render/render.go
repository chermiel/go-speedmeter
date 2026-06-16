package render

import (
	"fmt"
	"time"

	"speedmeter/internal/format"
)

const barWidth = 40

const (
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

type Data struct {
	Iface     string
	RxSpeed   float64
	TxSpeed   float64
	PeakRx    float64
	PeakTx    float64
	TotalRx   uint64
	TotalTx   uint64
	Runtime   time.Duration
}

func TUI(d Data, firstDraw bool) {
	w := barWidth - 2

	if firstDraw {
		fmt.Print("\033[H\033[2J")
		fmt.Println("╔══════════════════════════════════════════╗")
		fmt.Printf("║  🌐 SpeedMeter — %-19s ║\n", d.Iface)
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
		fmt.Print("\033[?25l")
	}

	fmt.Printf("\033[5;17H%s%-20s%s", colorGreen, format.Speed(d.RxSpeed), colorReset)
	fmt.Printf("\033[6;5H%s%s", format.Bar(d.RxSpeed/(d.PeakRx+1), w), colorReset)

	fmt.Printf("\033[8;16H%s%-20s%s", colorYellow, format.Speed(d.TxSpeed), colorReset)
	fmt.Printf("\033[9;5H%s%s", format.Bar(d.TxSpeed/(d.PeakTx+1), w), colorReset)

	fmt.Printf("\033[12;17H%-22s", format.Bytes(d.TotalRx))
	fmt.Printf("\033[13;17H%-22s", format.Bytes(d.TotalTx))
	fmt.Printf("\033[14;18H%-22s", format.Speed(d.PeakRx))
	fmt.Printf("\033[15;18H%-22s", format.Speed(d.PeakTx))
	fmt.Printf("\033[16;18H%-22s", d.Runtime)
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
