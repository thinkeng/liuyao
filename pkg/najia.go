package pkg

import (
	"fmt"
	"strings"
)

// 六爻纳甲映射表
var naJiaMap = map[string]map[string][]string{
	"乾": {"内卦": {"甲子", "甲寅", "甲辰"}, "外卦": {"壬午", "壬申", "壬戌"}},
	"坤": {"内卦": {"乙未", "乙巳", "乙卯"}, "外卦": {"癸丑", "癸亥", "癸酉"}},
	"震": {"内卦": {"庚子", "庚寅", "庚辰"}, "外卦": {"庚午", "庚申", "庚戌"}},
	"巽": {"内卦": {"辛丑", "辛亥", "辛酉"}, "外卦": {"辛未", "辛巳", "辛卯"}},
	"坎": {"内卦": {"戊寅", "戊辰", "戊午"}, "外卦": {"戊申", "戊戌", "戊子"}},
	"离": {"内卦": {"己卯", "己丑", "己亥"}, "外卦": {"己酉", "己未", "己巳"}},
	"艮": {"内卦": {"丙辰", "丙午", "丙申"}, "外卦": {"丙戌", "丙子", "丙寅"}},
	"兑": {"内卦": {"丁巳", "丁卯", "丁丑"}, "外卦": {"丁亥", "丁酉", "丁未"}},
}

// 八卦映射表 (三爻二进制 -> 卦名)
var trigramMap = map[string]string{
	// "111": "乾",
	// "000": "坤",
	// "001": "震",
	// "110": "巽",
	// "010": "坎",
	// "101": "离",
	// "100": "艮",
	// "011": "兑",
	"111": "乾",
	"000": "坤",
	"100": "震",
	"011": "巽",
	"010": "坎",
	"101": "离",
	"001": "艮",
	"110": "兑",
}

// 获取纳甲干支
func getNaJia(trigram, position string) []string {
	if trigramData, ok := naJiaMap[trigram]; ok {
		if ganzhi, ok := trigramData[position]; ok {
			return ganzhi
		}
	}
	return nil
}

// 解析六爻卦
func ParseHexagram(hexagramStr string) ([]string, error) {
	if len(hexagramStr) != 6 {
		return nil, fmt.Errorf("卦象长度必须为6位")
	}

	// 转换动爻表示法为阴阳爻
	normalized := strings.Map(func(r rune) rune {
		switch r {
		case '0', '1':
			return r
		default:
			return -1
		}
	}, hexagramStr)

	if len(normalized) != 6 {
		return nil, fmt.Errorf("包含无效字符，只允许0和1")
	}

	// 拆分内外卦
	innerTri := normalized[:3] // 内卦（下卦）
	outerTri := normalized[3:] // 外卦（上卦）

	// 获取卦名
	innerName, ok1 := trigramMap[innerTri]
	outerName, ok2 := trigramMap[outerTri]

	if !ok1 || !ok2 {
		return nil, fmt.Errorf("无效的卦象组合")
	}

	// 获取纳甲干支
	innerGanzhi := getNaJia(innerName, "内卦")
	outerGanzhi := getNaJia(outerName, "外卦")

	if innerGanzhi == nil || outerGanzhi == nil {
		return nil, fmt.Errorf("找不到对应的纳甲配置")
	}

	// 组合结果 (内卦+外卦)
	result := append(innerGanzhi, outerGanzhi...)
	return result, nil
}
