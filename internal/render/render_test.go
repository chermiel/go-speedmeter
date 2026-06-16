package render

import (
	"testing"
)

func TestTUINoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("TUI render panicked: %v", r)
		}
	}()

	d := Data{
		Iface:   "eth0",
		RxSpeed: 0,
		TxSpeed: 0,
		PeakRx:  0,
		PeakTx:  0,
		TotalRx: 0,
		TotalTx: 0,
	}

	TUI(d, true)
	Plain(d)
}
