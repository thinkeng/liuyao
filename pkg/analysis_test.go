package pkg

import (
	"testing"
)

func TestDetermineYongShen_NewCategories(t *testing.T) {
	tests := []struct {
		category string
		gender   string
		want     string
	}{
		{CategorySiblings, "Male", "兄弟"},
		{CategoryParents, "Male", "父母"},
		{CategoryChildren, "Male", "子孙"},
		{CategorySiblings, "Female", "兄弟"},
		// Regression checks
		{CategoryCareer, "Male", "官鬼"},
		{CategoryWealth, "Male", "妻财"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			if got := DetermineYongShen(tt.category, tt.gender); got != tt.want {
				t.Errorf("DetermineYongShen(%v, %v) = %v, want %v", tt.category, tt.gender, got, tt.want)
			}
		})
	}
}
