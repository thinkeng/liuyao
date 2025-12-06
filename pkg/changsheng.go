package pkg

// 地支顺序：子丑寅卯辰巳午未申酉戌亥
var dizhi = []string{
	"子", "丑", "寅", "卯", "辰", "巳",
	"午", "未", "申", "酉", "戌", "亥",
}

// 十二长生顺序
var changShengOrder = []string{
	"长生", "沐浴", "冠带", "临官",
	"帝旺", "衰", "病", "死",
	"墓", "绝", "胎", "养",
}

// 五行起长生表（核心，让长生计算正确！）
var wuXingChangShengStart = map[string]string{
	"木": "亥", // 木长生在亥
	"火": "寅", // 火长生在寅
	"土": "申", // 土长生在申（古法有申/巳两派，此处选申）
	"金": "巳", // 金长生在巳
	"水": "申", // 水长生在申（古法有申/酉两派，此处选申）
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

// GetChangSheng 计算某地支在某五行起长生体系下的长生位置
// wuxing: 木火土金水
// zhi: 要判断的地支，如 "午"
// 返回值: 长生十二神之一，如 "帝旺"
func GetChangSheng(wuxing, zhi string) string {
	// 找五行的起长生点地支
	startZhi, ok := wuXingChangShengStart[wuxing]
	if !ok {
		return ""
	}

	// 找起点索引
	startIdx := indexOf(dizhi, startZhi)
	// 找目标地支索引
	targetIdx := indexOf(dizhi, zhi)

	if startIdx == -1 || targetIdx == -1 {
		return ""
	}

	// 计算相对距离
	diff := (targetIdx - startIdx + 12) % 12

	return changShengOrder[diff]
}
