package pkg

import (
	"strings"
	"testing"
)

func TestDetermineYongShen_Gender(t *testing.T) {
	tests := []struct {
		name     string
		category string
		gender   string
		expected string
	}{
		{
			name:     "Marriage - Male",
			category: CategoryMarriage,
			gender:   "Male",
			expected: "妻财",
		},
		{
			name:     "Marriage - Female",
			category: CategoryMarriage,
			gender:   "Female",
			expected: "官鬼",
		},
		{
			name:     "Marriage - Unspecified",
			category: CategoryMarriage,
			gender:   "",
			expected: "妻财",
		},
		{
			name:     "Wealth - Male",
			category: CategoryWealth,
			gender:   "Male",
			expected: "妻财",
		},
		{
			name:     "Wealth - Female",
			category: CategoryWealth,
			gender:   "Female",
			expected: "妻财",
		},
		{
			name:     "Career - Male",
			category: CategoryCareer,
			gender:   "Male",
			expected: "官鬼",
		},
		{
			name:     "Career - Female",
			category: CategoryCareer,
			gender:   "Female",
			expected: "官鬼",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineYongShen(tt.category, tt.gender)
			if got != tt.expected {
				t.Errorf("DetermineYongShen(%s, %s) = %v, want %v", tt.category, tt.gender, got, tt.expected)
			}
		})
	}
}

func TestJudgeJiXiong_GenderMarriage(t *testing.T) {
	tests := []struct {
		name      string
		strength  string
		category  string
		gender    string
		expectJi  string
		containMs string
	}{
		{
			name:      "Female Marriage Strong",
			strength:  "强",
			category:  CategoryMarriage,
			gender:    "Female",
			expectJi:  "吉",
			containMs: "夫星得力",
		},
		{
			name:      "Female Marriage Weak",
			strength:  "弱",
			category:  CategoryMarriage,
			gender:    "Female",
			expectJi:  "凶",
			containMs: "提防感情冷淡",
		},
		{
			name:      "Male Marriage Strong",
			strength:  "强",
			category:  CategoryMarriage,
			gender:    "Male",
			expectJi:  "吉",
			containMs: "妻贤家富",
		},
		{
			name:      "Male Marriage Weak",
			strength:  "弱",
			category:  CategoryMarriage,
			gender:    "Male",
			expectJi:  "凶",
			containMs: "求财或感情不顺",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			judgment, details := JudgeJiXiong(tt.strength, tt.category, tt.gender)
			if judgment != tt.expectJi {
				t.Errorf("JudgeJiXiong(%s) = %v, want %v", tt.name, judgment, tt.expectJi)
			}
			found := false
			for _, d := range details {
				if tt.containMs != "" && (tt.containMs == d || (len(d) > len(tt.containMs) && d[:len(tt.containMs)] == tt.containMs) || (len(d) >= len(tt.containMs) && (d[len(d)-len(tt.containMs):] == tt.containMs || (len(d) > 10 && (d[5:5+len(tt.containMs)] == tt.containMs || (len(d) > 20 && (d[15:15+len(tt.containMs)] == tt.containMs))))))) {
					// Simple contain check since Chinese characters might be tricky with partials in a more complex way
					// Actually strings.Contains is better
				}
				if tt.containMs != "" && (strings.Contains(d, tt.containMs)) {
					found = true
					break
				}
			}
			if tt.containMs != "" && !found {
				t.Errorf("JudgeJiXiong(%s) details did not contain %q, details: %v", tt.name, tt.containMs, details)
			}
		})
	}
}
