package pkg

import (
	"fmt"
	"strings"
)

// GetShenSha Calculates Shen Sha for a given line's Earthly Branch based on Day Gan, Day Zhi, and Month Zhi.
func GetShenSha(dayGan, dayZhi, monthZhi, lineZhi string) string {
	var shenShaList []string

	// 1. Gui Ren (Nobleman) - Based on Day Gan
	// 甲戊并牛羊，乙己鼠猴乡，丙丁猪鸡位，壬癸兔蛇藏，六辛逢马虎，此是贵人方。
	// Jia/Wu -> Chou(Ox)/Wei(Goat)
	// Yi/Ji -> Zi(Rat)/Shen(Monkey)
	// Bing/Ding -> Hai(Pig)/You(Rooster)
	// Ren/Gui -> Mao(Rabbit)/Si(Snake)
	// Xin -> Wu(Horse)/Yin(Tiger)
	// Geng? Geng is usually associated with Jia/Wu or has its own.
	// Standard: 甲戊庚牛羊 (Jia/Wu/Geng -> Chou/Wei)
	if isGuiRen(dayGan, lineZhi) {
		shenShaList = append(shenShaList, "贵人")
	}

	// 2. Lu Shen (Prosperity) - Based on Day Gan
	// 甲禄在寅，乙禄在卯，丙戊禄在巳，丁己禄在午，庚禄在申，辛禄在酉，壬禄在亥，癸禄在子。
	if isLuShen(dayGan, lineZhi) {
		shenShaList = append(shenShaList, "禄神")
	}

	// 3. Yang Ren (Goat Blade) - Based on Day Gan
	// Before Lu is Yang Ren (usually).
	// 甲羊刃在卯，乙羊刃在辰(or Yin?), 丙戊羊刃在午，丁己羊刃在未(or Si?), 庚羊刃在酉，辛羊刃在戌(or Shen?), 壬羊刃在子，癸羊刃在丑(or Hai?)
	// Common formula: Yang Ren is the branch AFTER Lu (for Yang stems) or BEFORE Lu (for Yin stems)?
	// Or simplified:
	// 甲->卯, 乙->辰, 丙->午, 丁->未, 戊->午, 己->未, 庚->酉, 辛->戌, 壬->子, 癸->丑
	if isYangRen(dayGan, lineZhi) {
		shenShaList = append(shenShaList, "羊刃")
	}

	// 4. Wen Chang (Intelligence) - Based on Day Gan
	// 甲乙巳午报君知，丙戊申宫丁己鸡，庚猪辛鼠壬逢虎，癸人见兔入云梯。
	// 甲->巳, 乙->午, 丙->申, 丁->酉, 戊->申, 己->酉, 庚->亥, 辛->子, 壬->寅, 癸->卯
	if isWenChang(dayGan, lineZhi) {
		shenShaList = append(shenShaList, "文昌")
	}

	// 5. Yi Ma (Traveling Horse) - Based on Day Zhi (or Month Zhi, usually Day Zhi in Liu Yao)
	// 申子辰马在寅，寅午戌马在申，巳酉丑马在亥，亥卯未马在巳。
	if isYiMa(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "驿马")
	}

	// 6. Tao Hua (Peach Blossom) - Based on Day Zhi
	// 申子辰鸡叫乱人伦 (You)，寅午戌兔从卯里出 (Mao)，巳酉丑跃马南方走 (Wu)，亥卯未鼠子当头忌 (Zi)。
	if isTaoHua(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "桃花")
	}

	// 7. Jie Sha (Robbery Sha) - Based on Day Zhi
	// 申子辰在巳，寅午戌在亥，巳酉丑在寅，亥卯未在申。
	if isJieSha(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "劫煞")
	}

	// 8. Hua Gai (Elegant Cover) - Based on Day Zhi
	// 申子辰在辰，寅午戌在戌，巳酉丑在丑，亥卯未在未。
	if isHuaGai(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "华盖")
	}

	// 9. Jiang Xing (General Star) - Based on Day Zhi
	// 申子辰在子，寅午戌在午，巳酉丑在酉，亥卯未在卯。
	if isJiangXing(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "将星")
	}

	// 10. Mou Xing (Plotting Star) - Based on Day Zhi
	// 申子辰在戌，寅午戌在辰，巳酉丑在未，亥卯未在丑。
	if isMouXing(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "谋星")
	}

	// 11. Tian Xi (Heavenly Happiness) - Based on Month Zhi (Seasonal)
	// 正月戌, 二月亥, 三月子, 四月丑, 五月寅, 六月卯, 七月辰, 八月巳, 九月午, 十月未, 十一月申, 十二月酉
	// This is basically the opposite of the month branch? No.
	// Spring (Yin/Mao/Chen) -> Xu/Hai/Zi?
	// Let's use a map.
	if isTianXi(monthZhi, lineZhi) {
		shenShaList = append(shenShaList, "天喜")
	}

	// 12. Zai Sha (Calamity Sha) - Based on Day Zhi
	// 申子辰在午，寅午戌在子，巳酉丑在卯，亥卯未在酉。
	if isZaiSha(dayZhi, lineZhi) {
		shenShaList = append(shenShaList, "灾煞")
	}

	return strings.Join(shenShaList, " ")
}

// GetShenShaConfig Returns a summary of Shen Sha locations for the day.
func GetShenShaConfig(dayGan, dayZhi, monthZhi string) []string {
	var config []string
	earthlyBranches := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// Map to store Shen Sha -> Branches
	shenShaMap := make(map[string][]string)
	shenShaOrder := []string{
		"贵人", "禄神", "羊刃", "文昌", "驿马", "桃花", "劫煞", "华盖", "将星", "谋星", "天喜", "灾煞",
	}

	for _, zhi := range earthlyBranches {
		shaStr := GetShenSha(dayGan, dayZhi, monthZhi, zhi)
		if shaStr != "" {
			shas := strings.Split(shaStr, " ")
			for _, sha := range shas {
				shenShaMap[sha] = append(shenShaMap[sha], zhi)
			}
		}
	}

	for _, name := range shenShaOrder {
		if locations, ok := shenShaMap[name]; ok {
			config = append(config, fmt.Sprintf("%s:%s", name, strings.Join(locations, ",")))
		}
	}

	return config
}

// --- Helper Functions ---

func isGuiRen(gan, zhi string) bool {
	// 甲戊庚牛羊 (Chou, Wei)
	if (gan == "甲" || gan == "戊" || gan == "庚") && (zhi == "丑" || zhi == "未") {
		return true
	}
	// 乙己鼠猴乡 (Zi, Shen)
	if (gan == "乙" || gan == "己") && (zhi == "子" || zhi == "申") {
		return true
	}
	// 丙丁猪鸡位 (Hai, You)
	if (gan == "丙" || gan == "丁") && (zhi == "亥" || zhi == "酉") {
		return true
	}
	// 壬癸兔蛇藏 (Mao, Si)
	if (gan == "壬" || gan == "癸") && (zhi == "卯" || zhi == "巳") {
		return true
	}
	// 六辛逢马虎 (Wu, Yin)
	if gan == "辛" && (zhi == "午" || zhi == "寅") {
		return true
	}
	return false
}

func isLuShen(gan, zhi string) bool {
	mapLu := map[string]string{
		"甲": "寅", "乙": "卯", "丙": "巳", "丁": "午",
		"戊": "巳", "己": "午", "庚": "申", "辛": "酉",
		"壬": "亥", "癸": "子",
	}
	return mapLu[gan] == zhi
}

func isYangRen(gan, zhi string) bool {
	mapYangRen := map[string]string{
		"甲": "卯", "乙": "辰", "丙": "午", "丁": "未",
		"戊": "午", "己": "未", "庚": "酉", "辛": "戌",
		"壬": "子", "癸": "丑",
	}
	return mapYangRen[gan] == zhi
}

func isWenChang(gan, zhi string) bool {
	mapWenChang := map[string]string{
		"甲": "巳", "乙": "午", "丙": "申", "丁": "酉",
		"戊": "申", "己": "酉", "庚": "亥", "辛": "子",
		"壬": "寅", "癸": "卯",
	}
	return mapWenChang[gan] == zhi
}

func isYiMa(dayZhi, lineZhi string) bool {
	// 申子辰 -> 寅
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "寅" {
		return true
	}
	// 寅午戌 -> 申
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "申" {
		return true
	}
	// 巳酉丑 -> 亥
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "亥" {
		return true
	}
	// 亥卯未 -> 巳
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "巳" {
		return true
	}
	return false
}

func isTaoHua(dayZhi, lineZhi string) bool {
	// 申子辰 -> 酉
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "酉" {
		return true
	}
	// 寅午戌 -> 卯
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "卯" {
		return true
	}
	// 巳酉丑 -> 午
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "午" {
		return true
	}
	// 亥卯未 -> 子
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "子" {
		return true
	}
	return false
}

func isJieSha(dayZhi, lineZhi string) bool {
	// 申子辰 -> 巳
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "巳" {
		return true
	}
	// 寅午戌 -> 亥
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "亥" {
		return true
	}
	// 巳酉丑 -> 寅
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "寅" {
		return true
	}
	// 亥卯未 -> 申
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "申" {
		return true
	}
	return false
}

func isHuaGai(dayZhi, lineZhi string) bool {
	// 申子辰 -> 辰
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "辰" {
		return true
	}
	// 寅午戌 -> 戌
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "戌" {
		return true
	}
	// 巳酉丑 -> 丑
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "丑" {
		return true
	}
	// 亥卯未 -> 未
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "未" {
		return true
	}
	return false
}

func isJiangXing(dayZhi, lineZhi string) bool {
	// 申子辰 -> 子
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "子" {
		return true
	}
	// 寅午戌 -> 午
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "午" {
		return true
	}
	// 巳酉丑 -> 酉
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "酉" {
		return true
	}
	// 亥卯未 -> 卯
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "卯" {
		return true
	}
	return false
}

func isMouXing(dayZhi, lineZhi string) bool {
	// 申子辰 -> 戌
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "戌" {
		return true
	}
	// 寅午戌 -> 辰
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "辰" {
		return true
	}
	// 巳酉丑 -> 未
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "未" {
		return true
	}
	// 亥卯未 -> 丑
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "丑" {
		return true
	}
	return false
}

func isTianXi(monthZhi, lineZhi string) bool {
	// 正月(寅)戌, 二月(卯)亥, 三月(辰)子, 四月(巳)丑, 五月(午)寅, 六月(未)卯
	// 七月(申)辰, 八月(酉)巳, 九月(戌)午, 十月(亥)未, 十一月(子)申, 十二月(丑)酉
	mapTianXi := map[string]string{
		"寅": "戌", "卯": "亥", "辰": "子", "巳": "丑",
		"午": "寅", "未": "卯", "申": "辰", "酉": "巳",
		"戌": "午", "亥": "未", "子": "申", "丑": "酉",
	}
	return mapTianXi[monthZhi] == lineZhi
}

func isZaiSha(dayZhi, lineZhi string) bool {
	// 申子辰 -> 午
	if (dayZhi == "申" || dayZhi == "子" || dayZhi == "辰") && lineZhi == "午" {
		return true
	}
	// 寅午戌 -> 子
	if (dayZhi == "寅" || dayZhi == "午" || dayZhi == "戌") && lineZhi == "子" {
		return true
	}
	// 巳酉丑 -> 卯
	if (dayZhi == "巳" || dayZhi == "酉" || dayZhi == "丑") && lineZhi == "卯" {
		return true
	}
	// 亥卯未 -> 酉
	if (dayZhi == "亥" || dayZhi == "卯" || dayZhi == "未") && lineZhi == "酉" {
		return true
	}
	return false
}
