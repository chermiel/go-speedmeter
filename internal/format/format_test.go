package format

import "testing"

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
		ratio float64
		width int
		want  string
	}{
		{0, 10, "░░░░░░░░░░"},
		{1, 10, "██████████"},
		{0.5, 10, "█████░░░░░"},
		{0.33, 10, "███░░░░░░░"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Bar(tt.ratio, tt.width)
			if got != tt.want {
				t.Errorf("Bar(%v, %d) = %q, want %q", tt.ratio, tt.width, got, tt.want)
			}
		})
	}
}
