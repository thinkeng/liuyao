// Package pkg provides functionality for Liu Yao divination.
package pkg

import (
	"fmt"
	"strings"
)

// 三爻卦定义
type Trigram struct {
	name   string
	symbol string
	binary string
	wuXing string // 五行属性
}

type Gua struct {
	BenGua  []string
	BianGua []string
	Changed []bool
}

// 卦象信息
type GuaInfo struct {
	Position string // 爻位
	Ganzhi   string // 干支
	LiuShen  string // 六神
	YaoType  string // 爻类型
	LiuQin   string // 六亲
	ShiYing  string // 世应
	FuShen   string // 伏神 (如果有)
}

var (
	trigrams = []Trigram{
		{"乾", "☰", "111", "金"},
		{"兑", "☱", "110", "金"}, // Bottom-to-Top: 1,1,0
		{"离", "☲", "101", "火"},
		{"震", "☳", "100", "木"}, // Bottom-to-Top: 1,0,0
		{"巽", "☴", "011", "木"}, // Bottom-to-Top: 0,1,1
		{"坎", "☵", "010", "水"},
		{"艮", "☶", "001", "土"}, // Bottom-to-Top: 0,0,1
		{"坤", "☷", "000", "土"},
	}

	guaMap = map[string][2]int{
		// 乾宫 (0)
		"乾为天": {0, 0}, "天风姤": {0, 1}, "天山遁": {0, 2}, "天地否": {0, 3},
		"风地观": {0, 4}, "山地剥": {0, 5}, "火地晋": {0, 6}, "火天大有": {0, 7},

		// 兑宫 (1)
		"兑为泽": {1, 0}, "泽水困": {1, 1}, "泽地萃": {1, 2}, "泽山咸": {1, 3},
		"水山蹇": {1, 4}, "地山谦": {1, 5}, "雷山小过": {1, 6}, "雷泽归妹": {1, 7},

		// 离宫 (2)
		"离为火": {2, 0}, "火山旅": {2, 1}, "火风鼎": {2, 2}, "火水未济": {2, 3},
		"山水蒙": {2, 4}, "风水涣": {2, 5}, "天水讼": {2, 6}, "天火同人": {2, 7},

		// 震宫 (3)
		"震为雷": {3, 0}, "雷地豫": {3, 1}, "雷水解": {3, 2}, "雷风恒": {3, 3},
		"地风升": {3, 4}, "水风井": {3, 5}, "泽风大过": {3, 6}, "泽雷随": {3, 7},

		// 巽宫 (4)
		"巽为风": {4, 0}, "风天小畜": {4, 1}, "风火家人": {4, 2}, "风雷益": {4, 3},
		"天雷无妄": {4, 4}, "火雷噬嗑": {4, 5}, "山雷颐": {4, 6}, "山风蛊": {4, 7},

		// 坎宫 (5)
		"坎为水": {5, 0}, "水泽节": {5, 1}, "水雷屯": {5, 2}, "水火既济": {5, 3},
		"泽火革": {5, 4}, "雷火丰": {5, 5}, "地火明夷": {5, 6}, "地水师": {5, 7},

		// 艮宫 (6)
		"艮为山": {6, 0}, "山火贲": {6, 1}, "山天大畜": {6, 2}, "山泽损": {6, 3},
		"火泽睽": {6, 4}, "天泽履": {6, 5}, "风泽中孚": {6, 6}, "风山渐": {6, 7},

		// 坤宫 (7)
		"坤为地": {7, 0}, "地雷复": {7, 1}, "地泽临": {7, 2}, "地天泰": {7, 3},
		"雷天大壮": {7, 4}, "泽天夬": {7, 5}, "水天需": {7, 6}, "水地比": {7, 7},
	}
)

// 生成卦
func GenerateGua(inputs []string) Gua {
	benYao := make([]string, 6)
	changed := make([]bool, 6)
	binYao := make([]string, 6)

	for i, input := range inputs {
		count := strings.Count(input, "1")

		switch {
		case input == "111": // 老阳
			benYao[i] = "1"
			changed[i] = true
			binYao[i] = "0"
		case input == "000": // 老阴
			benYao[i] = "0"
			changed[i] = true
			binYao[i] = "1"
		case count == 2: // 少阴 (两阳一阴) - 2 Heads + 1 Tail = 8 (Yin)
			benYao[i] = "0"
			binYao[i] = "0"
		case count == 1: // 少阳 (一阳两阴) - 1 Head + 2 Tails = 7 (Yang)
			benYao[i] = "1"
			binYao[i] = "1"
		}
	}

	return Gua{
		BenGua:  benYao,
		Changed: changed,
		BianGua: binYao,
	}

}

// 解析投掷结果
func ParseToss(toss string) (string, string) {
	ones := 0
	for _, c := range toss {
		if c == '1' {
			ones++
		}
	}

	switch {
	case ones == 3:
		return "—○", "老阳"
	case ones == 0:
		return "⚋×", "老阴"
	case ones == 2:
		return "⚋ ", "少阴"
	case ones == 1:
		return "⚊ ", "少阳"
	default:
		return "??", "未知"
	}
}

// 确定卦名
func DetermineGuaName(binStr string) string {
	// 转换为二进制字符串 (从初爻到上爻)
	//binStr := strings.Join(yao, "")

	// 六十四卦映射
	guaMap := map[string]string{
		// 乾宫
		"111111": "乾为天", "011111": "天风姤", "001111": "天山遁", "000111": "天地否",
		"000011": "风地观", "000001": "山地剥", "000101": "火地晋", "111101": "火天大有",

		// 兑宫
		"110110": "兑为泽", "010110": "泽水困", "000110": "泽地萃", "001110": "泽山咸",
		"001010": "水山蹇", "001000": "地山谦", "001100": "雷山小过", "110100": "雷泽归妹",

		// 离宫
		"101101": "离为火", "001101": "火山旅", "011101": "火风鼎", "010101": "火水未济",
		"010001": "山水蒙", "010011": "风水涣", "010111": "天水讼", "101111": "天火同人",

		// 震宫
		"100100": "震为雷", "000100": "雷地豫", "010100": "雷水解", "011100": "雷风恒",
		"011000": "地风升", "011010": "水风井", "011110": "泽风大过", "100110": "泽雷随",

		// 巽宫
		"011011": "巽为风", "111011": "风天小畜", "101011": "风火家人", "100011": "风雷益",
		"100111": "天雷无妄", "100101": "火雷噬嗑", "100001": "山雷颐", "011001": "山风蛊",

		// 坎宫
		"010010": "坎为水", "110010": "水泽节", "100010": "水雷屯", "101010": "水火既济",
		"101110": "泽火革", "101100": "雷火丰", "101000": "地火明夷", "010000": "地水师",

		// 艮宫
		"001001": "艮为山", "101001": "山火贲", "111001": "山天大畜", "110001": "山泽损",
		"110101": "火泽睽", "110111": "天泽履", "110011": "风泽中孚", "001011": "风山渐",

		// 坤宫
		"000000": "坤为地", "100000": "地雷复", "110000": "地泽临", "111000": "地天泰",
		"111100": "雷天大壮", "111110": "泽天夬", "111010": "水天需", "000010": "水地比",
	}

	if name, ok := guaMap[binStr]; ok {
		return name
	}
	return "未知卦 (" + binStr + ")"
}

func GetGuaInfo(hexagramStr, dayGan string) ([]GuaInfo, error) {
	// 构建完整卦象信息
	result := make([]GuaInfo, 6)

	// 获取六神配置
	liuShen := LiuShenConfig[dayGan]

	ganzhiList, _ := ParseHexagram(hexagramStr)

	// 获取卦名和宫位信息
	guaName := DetermineGuaName(hexagramStr)
	palaceIndex, indexInPalace, _ := GetGuaPalace(guaName)
	palaceWuXing := GetPalaceWuXing(palaceIndex)

	// 获取世应位置
	shiPos, yingPos := GetShiYing(indexInPalace)

	// 1. Collect present Liu Qin
	presentLiuQin := make(map[string]bool)
	tempResults := make([]GuaInfo, 6)

	for i := 0; i < 6; i++ {
		// 获取爻类型
		yaoType, ok := yaoTypeMap[hexagramStr[i:i+1]]
		if !ok {
			yaoType = "未知爻"
		}

		// 计算六亲
		lineWuXing := GetWuXingFromGanZhi(ganzhiList[i])
		liuQin := GetLiuQin(palaceWuXing, lineWuXing)
		presentLiuQin[liuQin] = true

		// 计算世应
		shiYing := ""
		if i+1 == shiPos {
			shiYing = "世"
		} else if i+1 == yingPos {
			shiYing = "应"
		}

		tempResults[i] = GuaInfo{
			Position: yaoPositions[i],
			Ganzhi:   ganzhiList[i],
			LiuShen:  liuShen[i],
			YaoType:  yaoType,
			LiuQin:   liuQin,
			ShiYing:  shiYing,
		}
	}

	// 2. Check for missing Liu Qin (Fu Shen)
	// Standard 5 Liu Qin: 父母, 兄弟, 官鬼, 妻财, 子孙
	allLiuQin := []string{"父母", "兄弟", "官鬼", "妻财", "子孙"}
	missingLiuQin := make([]string, 0)
	for _, lq := range allLiuQin {
		if !presentLiuQin[lq] {
			missingLiuQin = append(missingLiuQin, lq)
		}
	}

	// 3. If there are missing Liu Qin, find them in Ben Gong Gua
	fuShenMap := make(map[int]string) // Index -> Fu Shen String
	if len(missingLiuQin) > 0 {
		benGongHex := GetBenGongGuaBinary(palaceIndex)
		if benGongHex != "" {
			benGongGanzhi, _ := ParseHexagram(benGongHex)
			// Ben Gong Gua belongs to the same palace, so Palace Wu Xing is the same
			for i, bgGanzhi := range benGongGanzhi {
				bgWuXing := GetWuXingFromGanZhi(bgGanzhi)
				bgLiuQin := GetLiuQin(palaceWuXing, bgWuXing)

				// Check if this Liu Qin is one of the missing ones
				for _, missing := range missingLiuQin {
					if bgLiuQin == missing {
						// Found a Fu Shen!
						// Format: "伏Name:Ganzhi" e.g. "伏父母:乙未"
						// Or just "父母:乙未"
						fuShenMap[i] = fmt.Sprintf("%s:%s", bgLiuQin, bgGanzhi)
					}
				}
			}
		}
	}

	// 4. Assign Fu Shen to results
	for i := 0; i < 6; i++ {
		result[i] = tempResults[i]
		if fs, ok := fuShenMap[i]; ok {
			result[i].FuShen = fs
		}
	}

	return result, nil
}

// GetBianGuaInfo 获取变卦信息 (六亲基于本卦宫位)
func GetBianGuaInfo(hexagramStr, dayGan, benGuaPalaceWuXing string) ([]GuaInfo, error) {
	// 基本信息与本卦相同
	result := make([]GuaInfo, 6)
	liuShen := LiuShenConfig[dayGan]
	ganzhiList, _ := ParseHexagram(hexagramStr)

	// 变卦的世应位置基于变卦本身
	guaName := DetermineGuaName(hexagramStr)
	_, indexInPalace, _ := GetGuaPalace(guaName)
	shiPos, yingPos := GetShiYing(indexInPalace)

	for i := 0; i < 6; i++ {
		yaoType, ok := yaoTypeMap[hexagramStr[i:i+1]]
		if !ok {
			yaoType = "未知爻"
		}

		// 计算六亲 (使用本卦宫位五行)
		lineWuXing := GetWuXingFromGanZhi(ganzhiList[i])
		liuQin := GetLiuQin(benGuaPalaceWuXing, lineWuXing)

		shiYing := ""
		if i+1 == shiPos {
			shiYing = "世"
		} else if i+1 == yingPos {
			shiYing = "应"
		}

		result[i] = GuaInfo{
			Position: yaoPositions[i],
			Ganzhi:   ganzhiList[i],
			LiuShen:  liuShen[i],
			YaoType:  yaoType,
			LiuQin:   liuQin,
			ShiYing:  shiYing,
		}
	}

	return result, nil
}

// GetGuaPalace 获取卦所在的宫位和宫内次序
func GetGuaPalace(name string) (int, int, bool) {
	if val, ok := guaMap[name]; ok {
		return val[0], val[1], true
	}
	return -1, -1, false
}

// GetPalaceWuXing 获取宫位的五行属性
func GetPalaceWuXing(palaceIndex int) string {
	if palaceIndex >= 0 && palaceIndex < len(trigrams) {
		return trigrams[palaceIndex].wuXing
	}
	return ""
}

// GetBenGongGuaBinary 获取本宫卦（首卦）的二进制字符串
func GetBenGongGuaBinary(palaceIndex int) string {
	// Palace Index: 0=Qian, 1=Dui, 2=Li, 3=Zhen, 4=Xun, 5=Kan, 6=Gen, 7=Kun
	// Ben Gong Gua is when Upper and Lower Trigrams are the same as the Palace Trigram.
	if palaceIndex >= 0 && palaceIndex < len(trigrams) {
		triBin := trigrams[palaceIndex].binary
		return triBin + triBin
	}
	return ""
}
