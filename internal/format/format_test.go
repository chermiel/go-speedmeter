package format

import (
	"testing"
	"unicode/utf8"
)

func TestSpeed(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{"zero", 0, "   0.00 bps"},
		{"negative", -100, "   0.00 bps"},
		{"bps", 500, " 500.00 bps"},
		{"Kbps", 1500, "   1.50 Kbps"},
		{"Mbps", 5_000_000, "   5.00 Mbps"},
		{"Gbps", 2_500_000_000, "   2.50 Gbps"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Speed(tt.input)
			if got != tt.want {
				t.Errorf("Speed(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		want  string
	}{
		{"zero", 0, "0 B"},
		{"B", 500, "500 B"},
		{"KB", 2048, "2.00 KB"},
		{"MB", 5_000_000, "4.77 MB"},
		{"GB", 5_000_000_000, "4.66 GB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bytes(tt.input)
			if got != tt.want {
				t.Errorf("Bytes(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBar(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		width int
		want  string
	}{
		{"zero", 0, 10, "░░░░░░░░░░"},
		{"full", 1, 10, "██████████"},
		{"half", 0.5, 10, "█████░░░░░"},
		{"third", 0.33, 10, "███░░░░░░░"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bar(tt.ratio, tt.width)
			if got != tt.want {
				t.Errorf("Bar(%v, %d) = %q, want %q", tt.ratio, tt.width, got, tt.want)
			}
		})
	}
}

func TestBarFractional(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		width int
		check func(t *testing.T, s string)
	}{
		{"zero", 0, 10, func(t *testing.T, s string) {
			if utf8.RuneCountInString(s) != 10 {
				t.Errorf("BarFractional(0, 10) length = %d, want 10", utf8.RuneCountInString(s))
			}
		}},
		{"full", 1, 10, func(t *testing.T, s string) {
			if s != "██████████" {
				t.Errorf("BarFractional(1, 10) = %q, want 10 full blocks", s)
			}
		}},
		{"half", 0.5, 10, func(t *testing.T, s string) {
			if utf8.RuneCountInString(s) != 10 {
				t.Errorf("BarFractional(0.5, 10) length = %d, want 10", utf8.RuneCountInString(s))
			}
		}},
		{"small_exact", 0.125, 5, func(t *testing.T, s string) {
			if utf8.RuneCountInString(s) != 5 {
				t.Errorf("length = %d, want 5", utf8.RuneCountInString(s))
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BarFractional(tt.ratio, tt.width)
			tt.check(t, got)
		})
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{0, "  0.0%"},
		{0.5, " 50.0%"},
		{1, "100.0%"},
		{-1, "  0.0%"},
		{2, "100.0%"},
	}
	for _, tt := range tests {
		got := Percent(tt.input)
		if got != tt.want {
			t.Errorf("Percent(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
