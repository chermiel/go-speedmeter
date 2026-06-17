package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"speedmeter/internal/render"
	"speedmeter/internal/stats"
	"speedmeter/internal/term"
)

func main() {
	iface := "eth0"
	if len(os.Args) > 1 {
		iface = os.Args[1]
	}

	useColor := term.IsTerminal()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	c := stats.NewCollector(iface)
	if err := c.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if useColor {
		fmt.Print(term.HideCursor())
		defer fmt.Print(term.ShowCursor())
	}

	startTime := time.Now()
	firstDraw := true
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Draw initial TUI skeleton immediately (no data yet, just layout)
	if useColor {
		d := render.Data{
			Iface:   iface,
			Runtime: 0,
		}
		render.TUI(d, true)
		firstDraw = false
	}

	for {
		select {
		case <-sigChan:
			return
		case <-ticker.C:
			rxSpeed, txSpeed, err := c.Tick()
			if err != nil {
				continue
			}

			runtime := time.Since(startTime).Round(time.Second)

			d := render.Data{
				Iface:   iface,
				RxSpeed: rxSpeed,
				TxSpeed: txSpeed,
				PeakRx:  c.PeakRx,
				PeakTx:  c.PeakTx,
				TotalRx: c.TotalRx,
				TotalTx: c.TotalTx,
				Runtime: runtime,
			}

			if useColor {
				render.TUI(d, firstDraw)
				firstDraw = false
			} else {
				render.Plain(d)
			}
		}
	}
}
