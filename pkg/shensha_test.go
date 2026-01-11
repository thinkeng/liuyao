package pkg

import (
	"testing"
)

func TestGetGuaShen(t *testing.T) {
	tests := []struct {
		name        string
		shiPosition int
		isYang      bool
		want        string
	}{
		// Yang Shi (Start Zi, Forward)
		{"Yang Shi - Line 1 (Zi)", 1, true, "子"},
		{"Yang Shi - Line 2 (Chou)", 2, true, "丑"},
		{"Yang Shi - Line 3 (Yin)", 3, true, "寅"},
		{"Yang Shi - Line 6 (Si)", 6, true, "巳"},

		// Yin Shi (Start Wu, Forward)
		{"Yin Shi - Line 1 (Wu)", 1, false, "午"},
		{"Yin Shi - Line 2 (Wei)", 2, false, "未"},
		{"Yin Shi - Line 3 (Shen)", 3, false, "申"},
		{"Yin Shi - Line 6 (亥)", 6, false, "亥"},

		// Boundary
		{"Values < 1", 0, true, ""},
		{"Values > 6", 7, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGuaShen(tt.shiPosition, tt.isYang); got != tt.want {
				t.Errorf("GetGuaShen(%d, %v) = %v, want %v", tt.shiPosition, tt.isYang, got, tt.want)
			}
		})
	}
}
