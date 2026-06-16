package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NetStat struct {
	RxBytes uint64
	TxBytes uint64
}

type Collector struct {
	Iface     string
	LastRx    uint64
	LastTx    uint64
	TotalRx   uint64
	TotalTx   uint64
	PeakRx    float64
	PeakTx    float64
}

func NewCollector(iface string) *Collector {
	return &Collector{Iface: iface}
}

func ReadNetStats(iface string) (NetStat, error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return NetStat{}, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines[2:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		name := strings.TrimRight(fields[0], ":")
		if name == iface {
			rx, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return NetStat{}, fmt.Errorf("parse rx: %w", err)
			}
			tx, err := strconv.ParseUint(fields[9], 10, 64)
			if err != nil {
				return NetStat{}, fmt.Errorf("parse tx: %w", err)
			}
			return NetStat{RxBytes: rx, TxBytes: tx}, nil
		}
	}
	return NetStat{}, fmt.Errorf("interface %s not found", iface)
}

func (c *Collector) Init() error {
	stat, err := ReadNetStats(c.Iface)
	if err != nil {
		return err
	}
	c.LastRx = stat.RxBytes
	c.LastTx = stat.TxBytes
	return nil
}

func (c *Collector) Tick() (rxSpeed, txSpeed float64, err error) {
	stat, err := ReadNetStats(c.Iface)
	if err != nil {
		return 0, 0, err
	}

	deltaRx := stat.RxBytes - c.LastRx
	deltaTx := stat.TxBytes - c.LastTx

	rxSpeed = float64(deltaRx) * 8
	txSpeed = float64(deltaTx) * 8

	c.TotalRx = stat.RxBytes
	c.TotalTx = stat.TxBytes

	if rxSpeed > c.PeakRx {
		c.PeakRx = rxSpeed
	}
	if txSpeed > c.PeakTx {
		c.PeakTx = txSpeed
	}

	c.LastRx = stat.RxBytes
	c.LastTx = stat.TxBytes

	return rxSpeed, txSpeed, nil
}
