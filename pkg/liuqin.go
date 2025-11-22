// Package pkg provides functionality for Liu Yao divination.
package pkg

// GetWuXingFromGanZhi returns the Wu Xing (Five Elements) of a GanZhi string.
// It looks at the second character (Earthly Branch).
func GetWuXingFromGanZhi(ganzhi string) string {
	if len(ganzhi) < 2 { // Should be at least 2 chars like "甲子"
		return ""
	}

	// Extract Earthly Branch (Zhi)
	// Assuming standard format like "甲子", "乙丑"
	// If it's just "子", "丑" etc, handle that too.
	zhi := ""
	runes := []rune(ganzhi)
	if len(runes) >= 2 {
		zhi = string(runes[1])
	} else {
		zhi = string(runes[0])
	}

	switch zhi {
	case "亥", "子":
		return "水"
	case "寅", "卯":
		return "木"
	case "巳", "午":
		return "火"
	case "申", "酉":
		return "金"
	case "辰", "戌", "丑", "未":
		return "土"
	default:
		return ""
	}
}

// GetLiuQin returns the Liu Qin (Six Relations) based on Palace Wu Xing and Line Wu Xing.
func GetLiuQin(palaceWuXing, lineWuXing string) string {
	if palaceWuXing == "" || lineWuXing == "" {
		return "未知"
	}

	// Relation Map: Palace (Subject) vs Line (Object)
	// Same -> Brothers (兄弟)
	// Palace generates Line -> Offspring (子孙)
	// Line generates Palace -> Parents (父母)
	// Palace controls Line -> Wealth (妻财)
	// Line controls Palace -> Official (官鬼)

	if palaceWuXing == lineWuXing {
		return "兄弟"
	}

	switch palaceWuXing {
	case "金":
		switch lineWuXing {
		case "水":
			return "子孙" // Metal generates Water
		case "木":
			return "妻财" // Metal controls Wood
		case "火":
			return "官鬼" // Fire controls Metal
		case "土":
			return "父母" // Earth generates Metal
		}
	case "木":
		switch lineWuXing {
		case "火":
			return "子孙" // Wood generates Fire
		case "土":
			return "妻财" // Wood controls Earth
		case "金":
			return "官鬼" // Metal controls Wood
		case "水":
			return "父母" // Water generates Wood
		}
	case "水":
		switch lineWuXing {
		case "木":
			return "子孙" // Water generates Wood
		case "火":
			return "妻财" // Water controls Fire
		case "土":
			return "官鬼" // Earth controls Water
		case "金":
			return "父母" // Metal generates Water
		}
	case "火":
		switch lineWuXing {
		case "土":
			return "子孙" // Fire generates Earth
		case "金":
			return "妻财" // Fire controls Metal
		case "水":
			return "官鬼" // Water controls Fire
		case "木":
			return "父母" // Wood generates Fire
		}
	case "土":
		switch lineWuXing {
		case "金":
			return "子孙" // Earth generates Metal
		case "水":
			return "妻财" // Earth controls Water
		case "木":
			return "官鬼" // Wood controls Earth
		case "火":
			return "父母" // Fire generates Earth
		}
	}

	return "未知"
}
