package pkg

import (
	"strings"
	"testing"
)

func TestGenerateGua_CoinTossLogic(t *testing.T) {
	// Test cases for Coin Toss interpretation
	// 1 = Head (Yang face), 0 = Tail (Yin face)

	// Case 1: 2 Heads, 1 Tail ("110") -> Shao Yin (Young Yin) -> Should be Yin line "0"
	// Current implementation likely returns "1" (Yang)
	input1 := []string{"110", "110", "110", "110", "110", "110"}
	gua1 := GenerateGua(input1)
	if gua1.BenGua[0] != "0" {
		t.Errorf("Expected '0' (Yin) for input '110' (2 Heads), got '%s'", gua1.BenGua[0])
	}

	// Case 2: 1 Head, 2 Tails ("100") -> Shao Yang (Young Yang) -> Should be Yang line "1"
	// Current implementation likely returns "0" (Yin)
	input2 := []string{"100", "100", "100", "100", "100", "100"}
	gua2 := GenerateGua(input2)
	if gua2.BenGua[0] != "1" {
		t.Errorf("Expected '1' (Yang) for input '100' (1 Head), got '%s'", gua2.BenGua[0])
	}

	// Case 3: 3 Heads ("111") -> Lao Yang (Old Yang) -> Should be Yang line "1" and Changing
	input3 := []string{"111", "111", "111", "111", "111", "111"}
	gua3 := GenerateGua(input3)
	if gua3.BenGua[0] != "1" {
		t.Errorf("Expected '1' (Yang) for input '111' (3 Heads), got '%s'", gua3.BenGua[0])
	}
	if !gua3.Changed[0] {
		t.Errorf("Expected Changed=true for input '111'")
	}
	if gua3.BianGua[0] != "0" {
		t.Errorf("Expected BianGua '0' (Yin) for input '111', got '%s'", gua3.BianGua[0])
	}

	// Case 4: 0 Heads ("000") -> Lao Yin (Old Yin) -> Should be Yin line "0" and Changing
	input4 := []string{"000", "000", "000", "000", "000", "000"}
	gua4 := GenerateGua(input4)
	if gua4.BenGua[0] != "0" {
		t.Errorf("Expected '0' (Yin) for input '000' (0 Heads), got '%s'", gua4.BenGua[0])
	}
	if !gua4.Changed[0] {
		t.Errorf("Expected Changed=true for input '000'")
	}
	if gua4.BianGua[0] != "1" {
		t.Errorf("Expected BianGua '1' (Yang) for input '000', got '%s'", gua4.BianGua[0])
	}
}

func TestGetShiYing(t *testing.T) {
	tests := []struct {
		index int
		shi   int
		ying  int
	}{
		{0, 6, 3}, // Ben Gua
		{1, 1, 4}, // 1st Change
		{2, 2, 5}, // 2nd Change
		{3, 3, 6}, // 3rd Change
		{4, 4, 1}, // 4th Change
		{5, 5, 2}, // 5th Change
		{6, 4, 1}, // You Hun
		{7, 3, 6}, // Gui Hun
	}

	for _, tt := range tests {
		s, y := GetShiYing(tt.index)
		if s != tt.shi || y != tt.ying {
			t.Errorf("GetShiYing(%d) = (%d, %d), want (%d, %d)", tt.index, s, y, tt.shi, tt.ying)
		}
	}
}

func TestGetLiuQin(t *testing.T) {
	tests := []struct {
		palace string
		line   string
		want   string
	}{
		{"金", "金", "兄弟"},
		{"金", "水", "子孙"},
		{"金", "木", "妻财"},
		{"金", "火", "官鬼"},
		{"金", "土", "父母"},
		{"火", "土", "子孙"}, // Fire generates Earth
		{"火", "木", "父母"}, // Wood generates Fire
	}

	for _, tt := range tests {
		got := GetLiuQin(tt.palace, tt.line)
		if got != tt.want {
			t.Errorf("GetLiuQin(%s, %s) = %s, want %s", tt.palace, tt.line, got, tt.want)
		}
	}
}

func TestGetGuaInfo_MengGua(t *testing.T) {
	// Meng Gua (Mountain Water Meng) - 010001
	// Palace: Li (Fire)
	// Shi: 4, Ying: 1
	hexagram := "010001"
	dayGan := "庚" // Day Stem (affects Liu Shen)

	infos, err := GetGuaInfo(hexagram, dayGan)
	if err != nil {
		t.Fatalf("GetGuaInfo failed: %v", err)
	}

	// Check Shi/Ying
	if infos[3].ShiYing != "世" { // 4th line (index 3)
		t.Errorf("Expected 4th line to be Shi, got %s", infos[3].ShiYing)
	}
	if infos[0].ShiYing != "应" { // 1st line (index 0)
		t.Errorf("Expected 1st line to be Ying, got %s", infos[0].ShiYing)
	}

	// Check Liu Qin
	// Palace is Fire.
	// Line 1 (Chu): Kan Inner -> Wu Yin (Wood). Wood generates Fire -> Parents (父母).
	if infos[0].LiuQin != "父母" {
		t.Errorf("Expected 1st line LiuQin to be 父母, got %s (Ganzhi: %s)", infos[0].LiuQin, infos[0].Ganzhi)
	}

	// Line 4 (Four): Gen Outer -> Bing Xu (Earth). Fire generates Earth -> Offspring (子孙).
	if infos[3].LiuQin != "子孙" {
		t.Errorf("Expected 4th line LiuQin to be 子孙, got %s (Ganzhi: %s)", infos[3].LiuQin, infos[3].Ganzhi)
	}
}

func TestBianGua_LiuQin_Logic(t *testing.T) {
	// Case: Qian (Metal) changing to Kun (Earth)
	// Ben Gua: Qian (111111) -> Palace: Metal
	// Bian Gua: Kun (000000) -> Palace: Earth (Normally)
	// But for Bian Gua Liu Qin, we must use Ben Gua's Palace (Metal).

	// Kun Line 1: Yi Wei (Earth)
	// If using Kun Palace (Earth): Earth vs Earth = Brother
	// If using Qian Palace (Metal): Metal vs Earth (Earth generates Metal) = Parent

	benGuaPalaceWuXing := "金" // Qian
	bianHexagram := "000000"  // Kun
	dayGan := "甲"             // Unused

	// Test GetBianGuaInfo
	infos, err := GetBianGuaInfo(bianHexagram, dayGan, benGuaPalaceWuXing)
	if err != nil {
		t.Fatalf("GetBianGuaInfo failed: %v", err)
	}

	// Kun Line 1: Yi Wei (Earth)
	// Palace: Metal (from Ben Gua)
	// Earth generates Metal -> Parent (父母)
	if infos[0].LiuQin != "父母" {
		t.Errorf("Expected Bian Gua Line 1 LiuQin to be 父母, got %s (Ganzhi: %s)", infos[0].LiuQin, infos[0].Ganzhi)
	}
}

func TestFuShen_Logic(t *testing.T) {
	// Case: Gou Gua (Heaven Wind Gou) - 011111
	// Palace: Qian (Metal)
	// Lines:
	// 1: Chou Earth (Parent)
	// 2: Hai Water (Offspring)
	// 3: You Metal (Brother)
	// 4: Wu Fire (Official)
	// 5: Shen Metal (Brother)
	// 6: Xu Earth (Parent)
	// Missing: Wealth (Wood)

	// Ben Gong Gua: Qian (111111)
	// Line 2: Yin Wood (Wealth) -> Fu Shen for Line 2?
	// Wait, Fu Shen rules:
	// Look at Ben Gong Gua.
	// Qian Ben Gong:
	// 1: Zi Water (Offspring)
	// 2: Yin Wood (Wealth)  <-- This is the missing Wealth
	// 3: Chen Earth (Parent)
	// 4: Wu Fire (Official)
	// 5: Shen Metal (Brother)
	// 6: Xu Earth (Parent)

	// So Wealth (Wood) is at Line 2 in Ben Gong Gua.
	// It should appear as Fu Shen under Line 2 of Gou Gua.

	hexagram := "011111" // Gou Gua
	dayGan := "甲"

	infos, err := GetGuaInfo(hexagram, dayGan)
	if err != nil {
		t.Fatalf("GetGuaInfo failed: %v", err)
	}

	// Check Line 2 (Index 1) for Fu Shen
	if !strings.Contains(infos[1].FuShen, "妻财") {
		t.Errorf("Expected Fu Shen 'Wealth' (妻财) at Line 2, got '%s'", infos[1].FuShen)
	}
}
