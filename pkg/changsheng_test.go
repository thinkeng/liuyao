package pkg

import (
	"testing"
)

func TestGetChangSheng(t *testing.T) {
	tests := []struct {
		name     string
		wuxing   string
		zhi      string
		expected string
	}{
		// 木长生在亥
		{"Wood - Hai (ChangSheng)", "木", "亥", "长生"},
		{"Wood - Zi (MuYu)", "木", "子", "沐浴"},
		{"Wood - Mao (YiWang)", "木", "卯", "帝旺"},
		{"Wood - Wei (Mu)", "木", "未", "墓"},

		// 火长生在寅
		{"Fire - Yin (ChangSheng)", "火", "寅", "长生"},
		{"Fire - Wu (YiWang)", "火", "午", "帝旺"},
		{"Fire - Xu (Mu)", "火", "戌", "墓"},

		// 金长生在巳
		{"Metal - Si (ChangSheng)", "金", "巳", "长生"},
		{"Metal - You (YiWang)", "金", "酉", "帝旺"},
		{"Metal - Chou (Mu)", "金", "丑", "墓"},

		// 水/土长生在申 (水土同宫)
		{"Water - Shen (ChangSheng)", "水", "申", "长生"},
		{"Water - Zi (YiWang)", "水", "子", "帝旺"},
		{"Water - Chen (Mu)", "水", "辰", "墓"},

		{"Earth - Shen (ChangSheng)", "土", "申", "长生"},
		{"Earth - Zi (YiWang)", "土", "子", "帝旺"},
		{"Earth - Chen (Mu)", "土", "辰", "墓"},

		// Invalid cases
		{"Invalid WuXing", "Unknown", "子", ""},
		{"Invalid Zhi", "木", "Unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetChangSheng(tt.wuxing, tt.zhi)
			if got != tt.expected {
				t.Errorf("GetChangSheng(%q, %q) = %q, want %q", tt.wuxing, tt.zhi, got, tt.expected)
			}
		})
	}
}
