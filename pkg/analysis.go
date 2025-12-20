package pkg

import (
	"fmt"
	"strings"
	"time"
)

// Question Category Constants
const (
	CategoryCareer   = "Career"   // 求官/工作
	CategoryWealth   = "Wealth"   // 求财
	CategoryMarriage = "Marriage" // 婚姻
	CategoryStudy    = "Study"    // 学业/文书
	CategorySafety   = "Safety"   // 平安/子嗣
	CategoryHealth   = "Health"   // 健康 (Self)
)

// AnalysisContext holds the input data for analysis
type AnalysisContext struct {
	GuaHexagram  string    // Binary string of the Ben Gua (e.g., "111000")
	BianHexagram string    // Binary string of the Bian Gua (e.g., "111111")
	Changed      []bool    // Array indicating which lines are moving
	DayGan       string    // Day Heavenly Stem
	DayZhi       string    // Day Earthly Branch
	MonthZhi     string    // Month Earthly Branch
	DayXunKong   string    // Day Xun Kong (Empty Branches)
	Category     string    // Question Category
	Gender       string    // Gender: "Male" or "Female"
	Date         time.Time // Date of divination
}

// AnalysisResult holds the output of the analysis
type AnalysisResult struct {
	YongShen      string    // The Use God (e.g., "官鬼")
	YongShenYao   GuaInfo   // The specific Yao representing the Use God
	YongShenIndex int       // Index of the Use God Yao (0-5)
	Strength      string    // Overall strength description
	Judgment      string    // "Ji" (Auspicous) or "Xiong" (Inauspicious)
	Details       []string  // Detailed analysis steps
	GuaName       string    // 卦名
	GuaCi         string    // 卦辞
	CoreMeaning   string    // 核心意象
	MovingYaos    []YaoText // 动爻文本信息
}

// Analyze performs the hexagram analysis
func Analyze(ctx AnalysisContext) (AnalysisResult, error) {
	result := AnalysisResult{
		Details:    make([]string, 0),
		MovingYaos: make([]YaoText, 0),
	}

	// 1. Determine Use God (Yong Shen)
	yongShen := DetermineYongShen(ctx.Category, ctx.Gender)
	result.YongShen = yongShen

	// --- Enrich Text Info (Gua & Yao) ---
	guaName := DetermineGuaName(ctx.GuaHexagram)
	result.GuaName = guaName

	// Query Gua Text
	guaText, _, errGua := QueryGuaAndYaoCi(guaName, "")
	if errGua == nil {
		result.GuaCi = guaText.GuaCi
		result.CoreMeaning = guaText.CoreMeaning
	}

	// Query Moving Yao Text
	for i, changed := range ctx.Changed {
		if changed {
			// Construct Yao Name (e.g., "初九", "六二")
			// Binary string index 0 is Bottom (Line 1)
			bit := string(ctx.GuaHexagram[i])
			yaoName := GetYaoName(i, bit)

			_, yaoText, errYao := QueryGuaAndYaoCi(guaName, yaoName)
			if errYao == nil {
				result.MovingYaos = append(result.MovingYaos, yaoText)
			}
		}
	}
	// ------------------------------------

	categoryCn := ctx.Category
	switch ctx.Category {
	case CategoryWealth:
		categoryCn = "求财"
	case CategoryCareer:
		categoryCn = "求官/工作"
	case CategoryMarriage:
		categoryCn = "婚姻"
	case CategoryStudy:
		categoryCn = "学业"
	case CategorySafety:
		categoryCn = "平安"
	case CategoryHealth:
		categoryCn = "健康"
	}

	result.Details = append(result.Details, fmt.Sprintf("求测事项: %s -> 用神: %s", categoryCn, yongShen))

	// 2. Get Gua Info to find the Use God Yao
	guaInfo, err := GetGuaInfo(ctx.GuaHexagram, ctx.DayGan)
	if err != nil {
		return result, err
	}

	// Find the Use God in the hexagram
	// Priority:
	// 1. Dong Yao (Moving Line) with the right Liu Qin? (Usually we look for the specific Liu Qin first)
	// 2. If multiple, usually pick the one that is Moving, or if none moving, the one with better position/strength.
	// For simplicity V1: Pick the first occurrence, or if multiple, prefer the Moving one.
	// Actually, standard rule:
	// - If one appears, take it.
	// - If multiple appear:
	//   - Pick the Moving one.
	//   - If both moving or both static, pick the one closer to Shi Yao or with better strength.
	//   - For now, let's pick the first one found, but prioritize Moving if possible?
	//   Let's just pick the first one for now and refine later.

	foundIndex := -1

	// Search for Yong Shen
	candidates := []int{}
	for i, info := range guaInfo {
		if info.LiuQin == yongShen {
			candidates = append(candidates, i)
		}
	}

	if len(candidates) == 0 {
		// Not found in Ben Gua -> Check Fu Shen (Hidden Spirit)
		for i, info := range guaInfo {
			if strings.Contains(info.FuShen, yongShen) {
				foundIndex = i
				result.Details = append(result.Details, fmt.Sprintf("用神 %s 不现，伏藏于 %s 下 (爻位: %s)", yongShen, info.LiuQin, info.Position))
				break
			}
		}

		if foundIndex == -1 {
			return result, fmt.Errorf("用神 %s 不现且未伏藏", yongShen)
		}
	} else if len(candidates) == 1 {
		// Only one candidate, use it
		foundIndex = candidates[0]
	} else {
		// Multiple candidates - apply priority rules
		// Priority 1: 持世 (Holds World)
		for _, idx := range candidates {
			if guaInfo[idx].ShiYing == "世" {
				foundIndex = idx
				result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取持世之爻 (爻位: %s)", yongShen, guaInfo[idx].Position))
				break
			}
		}

		// Priority 2: 发动 (Moving Line)
		if foundIndex == -1 && len(ctx.Changed) > 0 {
			for _, idx := range candidates {
				if len(ctx.Changed) > idx && ctx.Changed[idx] {
					foundIndex = idx
					result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取发动之爻 (爻位: %s)", yongShen, guaInfo[idx].Position))
					break
				}
			}
		}

		// Priority 3: 临日月 (Aligns with Day or Month)
		if foundIndex == -1 {
			for _, idx := range candidates {
				yaoZhi := string([]rune(guaInfo[idx].Ganzhi)[1])
				if yaoZhi == ctx.MonthZhi {
					foundIndex = idx
					result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取临月建之爻 (爻位: %s)", yongShen, guaInfo[idx].Position))
					break
				}
			}
		}

		if foundIndex == -1 {
			for _, idx := range candidates {
				yaoZhi := string([]rune(guaInfo[idx].Ganzhi)[1])
				if yaoZhi == ctx.DayZhi {
					foundIndex = idx
					result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取临日辰之爻 (爻位: %s)", yongShen, guaInfo[idx].Position))
					break
				}
			}
		}

		// Priority 4: 临应爻 (Aligns with Response)
		if foundIndex == -1 {
			for _, idx := range candidates {
				if guaInfo[idx].ShiYing == "应" {
					foundIndex = idx
					result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取临应爻 (爻位: %s)", yongShen, guaInfo[idx].Position))
					break
				}
			}
		}

		// Fallback: Take the first one
		if foundIndex == -1 {
			foundIndex = candidates[0]
			result.Details = append(result.Details, fmt.Sprintf("出现多个 %s，取第一个 (爻位: %s)", yongShen, guaInfo[foundIndex].Position))
		}
	}

	result.YongShenIndex = foundIndex
	result.YongShenYao = guaInfo[foundIndex]

	// Phase 1.5: Fu Shen (Hidden Spirit) Handling
	isFuShen := result.YongShenYao.FuShen != "" && strings.Contains(result.YongShenYao.FuShen, yongShen)
	var fuShenGanzhi string
	fuShenScore := 0

	if isFuShen {
		fuShenParts := strings.Split(result.YongShenYao.FuShen, ":")
		if len(fuShenParts) == 2 {
			fuShenGanzhi = fuShenParts[1]
			fuShenWuXing := GetWuXingFromGanZhi(fuShenGanzhi)
			fuShenZhi := string([]rune(fuShenGanzhi)[1])
			feiShenWuXing := GetWuXingFromGanZhi(result.YongShenYao.Ganzhi)
			monthWuXing := GetWuXing(ctx.MonthZhi)
			dayWuXing := GetWuXing(ctx.DayZhi)

			// 1. Month/Day Influence (Decisive) - Added to details first
			monthStr := GetMonthStrength(fuShenWuXing, monthWuXing)
			dayStr := GetDayStrength(fuShenWuXing, dayWuXing)
			result.Details = append(result.Details, fmt.Sprintf("伏神月日影响: 月建(%s) %s, 日辰(%s) %s", monthWuXing, monthStr, dayWuXing, dayStr))

			// 2. Fei-Fu Relationship (Supportive)
			feiToFuRelation := GetRelation(feiShenWuXing, fuShenWuXing)
			fuToFeiRelation := GetRelation(fuShenWuXing, feiShenWuXing)

			fuShenDetail := fmt.Sprintf("飞伏关系: 飞神 %s (%s) ", result.YongShenYao.LiuQin, result.YongShenYao.Ganzhi)
			if feiToFuRelation == "Sheng" {
				fuShenDetail += "生 伏神 -> 飞生伏 (吉)"
				fuShenScore += 2
			} else if feiToFuRelation == "Ke" {
				fuShenDetail += "克 伏神 -> 飞克伏 (凶)"
				fuShenScore -= 2
			} else if fuToFeiRelation == "Sheng" {
				fuShenDetail += "被 伏神 生 -> 伏生飞 (泄气)"
				fuShenScore -= 1
			} else if fuToFeiRelation == "Ke" {
				fuShenDetail += "被 伏神 克 -> 伏克飞 (吉)"
				fuShenScore += 1
			} else {
				fuShenDetail += "与 伏神 比和/无生克"
			}
			result.Details = append(result.Details, fuShenDetail)

			// 3. Special States
			if CheckXunKong(fuShenGanzhi, ctx.DayXunKong) {
				result.Details = append(result.Details, "伏神特殊状态: 旬空")
				fuShenScore -= 2
			}
			if IsChong(ctx.MonthZhi, fuShenZhi) {
				result.Details = append(result.Details, "伏神特殊状态: 月破")
				fuShenScore -= 4
			}
		}
	}

	// Phase 2: Strength Analysis
	// Need Bian Gua Info for the Use God Line
	var bianYaoInfo *GuaInfo
	isMoving := false
	if len(ctx.Changed) > result.YongShenIndex && ctx.Changed[result.YongShenIndex] {
		isMoving = true
		// Get Bian Gua Info
		// We need to call GetBianGuaInfo. But wait, we don't have it handy here unless we call it.
		// Or we can just get it for the specific line if we have the Bian Hexagram.
		// Let's call GetBianGuaInfo for the whole hexagram to be safe/consistent.
		// We need Ben Gua Palace Wu Xing.
		guaName := DetermineGuaName(ctx.GuaHexagram)
		palaceIndex, _, _ := GetGuaPalace(guaName)
		palaceWuXing := GetPalaceWuXing(palaceIndex)

		bianResult, err := GetBianGuaInfo(ctx.BianHexagram, ctx.DayGan, palaceWuXing)
		if err == nil && len(bianResult) > result.YongShenIndex {
			bianYaoInfo = &bianResult[result.YongShenIndex]
		}
	}

	// Use virtual YaoInfo if Fu Shen
	targetYaoInfo := result.YongShenYao
	if isFuShen && fuShenGanzhi != "" {
		targetYaoInfo.Ganzhi = fuShenGanzhi
		targetYaoInfo.FuShen = "" // Don't trigger recursive Fu Shen checks in CalculateStrength
	}

	strength, strengthDetails := CalculateStrength(targetYaoInfo, bianYaoInfo, isMoving, ctx.MonthZhi, ctx.DayZhi, ctx.DayXunKong)

	// Apply Fu Shen score adjustment if applicable
	if isFuShen {
		if fuShenScore >= 2 {
			if strength == "弱" {
				strength = "中平 (飞神生助)"
			} else if strength == "中平" {
				strength = "强"
			}
		} else if fuShenScore <= -2 {
			if strength == "强" {
				strength = "中平 (飞神克制)"
			} else if strength == "中平" {
				strength = "弱"
			}
		}
	}

	result.Strength = strength
	result.Details = append(result.Details, strengthDetails...)

	// Comprehensive Line-by-Line Analysis
	result.Details = append(result.Details, "\n【各爻详细分析】")

	// Get Bian Gua Info for all lines if there are moving lines
	var bianGuaInfoAll []GuaInfo
	if len(ctx.Changed) > 0 {
		guaName := DetermineGuaName(ctx.GuaHexagram)
		palaceIndex, _, _ := GetGuaPalace(guaName)
		palaceWuXing := GetPalaceWuXing(palaceIndex)
		bianResult, err := GetBianGuaInfo(ctx.BianHexagram, ctx.DayGan, palaceWuXing)
		if err == nil {
			bianGuaInfoAll = bianResult
		}
	}

	monthWuXing := GetWuXing(ctx.MonthZhi)
	dayWuXing := GetWuXing(ctx.DayZhi)

	// Analyze each line from top to bottom
	for i := len(guaInfo) - 1; i >= 0; i-- {
		lineInfo := guaInfo[i]
		lineWuXing := GetWuXingFromGanZhi(lineInfo.Ganzhi)
		lineZhi := string([]rune(lineInfo.Ganzhi)[1])

		lineDetail := fmt.Sprintf("%s:", lineInfo.Position)

		// Month relationship
		if ctx.MonthZhi == lineZhi {
			lineDetail += "值月建"
		} else if IsSheng(monthWuXing, lineWuXing) {
			lineDetail += "月建生"
		} else if IsKe(monthWuXing, lineWuXing) {
			lineDetail += "月建克"
		} else {
			monthStrength := GetMonthStrength(lineWuXing, monthWuXing)
			lineDetail += fmt.Sprintf("月上%s", monthStrength)
		}

		// Day relationship
		if ctx.DayZhi == lineZhi {
			lineDetail += " 临日辰"
		} else if IsSheng(dayWuXing, lineWuXing) {
			lineDetail += " 日辰生"
		} else if IsKe(dayWuXing, lineWuXing) {
			lineDetail += " 日辰克"
		} else if IsChong(ctx.DayZhi, lineZhi) {
			lineDetail += " 日辰冲"
		}

		// Check for 日破
		if IsChong(ctx.DayZhi, lineZhi) {
			monthStrength := GetMonthStrength(lineWuXing, monthWuXing)
			if !IsStrong(monthStrength) {
				lineDetail += " 日破"
			}
		}

		// Advanced Relationships with Day/Month
		if he := CheckLiuHe(lineZhi, ctx.MonthZhi); he != "" {
			lineDetail += fmt.Sprintf(" 月合(%s)", he)
		}
		if he := CheckLiuHe(lineZhi, ctx.DayZhi); he != "" {
			lineDetail += fmt.Sprintf(" 六合(%s)", he)
		}

		if xing := CheckXing(lineZhi, ctx.DayZhi); xing != "" {
			lineDetail += fmt.Sprintf(" 日%s", xing)
		}
		if CheckLiuHai(lineZhi, ctx.DayZhi) {
			lineDetail += " 日害"
		}

		// Check for moving line
		if len(ctx.Changed) > i && ctx.Changed[i] {
			if len(bianGuaInfoAll) > i {
				bianLineInfo := bianGuaInfoAll[i]
				bianLineWuXing := GetWuXingFromGanZhi(bianLineInfo.Ganzhi)

				// Check moving -> changed relationship (动爻对变爻的关系)
				movingToChangedRelation := GetRelation(lineWuXing, bianLineWuXing)
				// Check changed -> moving relationship (变爻对动爻的关系)
				changedToMovingRelation := GetRelation(bianLineWuXing, lineWuXing)

				if changedToMovingRelation == "Sheng" {
					lineDetail += " 化回头生(变生动)"
				} else if changedToMovingRelation == "Ke" {
					lineDetail += " 化回头克(变克动)"
				} else if movingToChangedRelation == "Sheng" {
					lineDetail += " 化泄气(动生变)"
				} else if movingToChangedRelation == "Ke" {
					lineDetail += " 化克(动克变)"
				}

				jinTui := CheckJinTui(lineInfo.Ganzhi, bianLineInfo.Ganzhi)
				if jinTui == "Jin Shen" {
					lineDetail += " 化进神"
				} else if jinTui == "Tui Shen" {
					lineDetail += " 化退神"
				}
			}
		}

		// Check interactions with other lines
		interactionCount := 0
		for j := len(guaInfo) - 1; j >= 0; j-- {
			if i == j {
				continue
			}
			otherInfo := guaInfo[j]
			otherWuXing := GetWuXingFromGanZhi(otherInfo.Ganzhi)
			otherZhi := string([]rune(otherInfo.Ganzhi)[1])

			// Check Chong
			if IsChong(lineZhi, otherZhi) {
				lineDetail += fmt.Sprintf(" 冲%s", otherInfo.Position)
				interactionCount++
			}

			// Check Sheng/Ke from other moving lines
			if len(ctx.Changed) > j && ctx.Changed[j] {
				relation := GetRelation(otherWuXing, lineWuXing)
				if relation == "Sheng" {
					lineDetail += fmt.Sprintf(" %s生", otherInfo.Position)
					interactionCount++
				} else if relation == "Ke" {
					lineDetail += fmt.Sprintf(" %s克", otherInfo.Position)
					interactionCount++
				}

				// Check Advanced Relations with Moving Lines
				if he := CheckLiuHe(lineZhi, otherZhi); he != "" {
					lineDetail += fmt.Sprintf(" %s合(%s)", otherInfo.Position, he)
				}
				if CheckLiuHai(lineZhi, otherZhi) {
					lineDetail += fmt.Sprintf(" %s害", otherInfo.Position)
				}
				if xing := CheckXing(lineZhi, otherZhi); xing != "" {
					lineDetail += fmt.Sprintf(" %s%s", otherInfo.Position, xing)
				}
			}
		}

		// Add changed line info for moving lines
		if len(ctx.Changed) > i && ctx.Changed[i] && len(bianGuaInfoAll) > i {
			bianLineInfo := bianGuaInfoAll[i]
			bianLineWuXing := GetWuXingFromGanZhi(bianLineInfo.Ganzhi)
			bianLineZhi := string([]rune(bianLineInfo.Ganzhi)[1])

			lineDetail += fmt.Sprintf("\n  变爻→%s %s", bianLineInfo.LiuQin, bianLineInfo.Ganzhi)

			if ctx.MonthZhi == bianLineZhi {
				lineDetail += " 值月建"
			} else if IsSheng(monthWuXing, bianLineWuXing) {
				lineDetail += " 月建生"
			}

			if ctx.DayZhi == bianLineZhi {
				lineDetail += " 临日辰"
			} else if IsSheng(dayWuXing, bianLineWuXing) {
				lineDetail += " 日辰生"
			}
		}

		// Add hidden spirit info
		if lineInfo.FuShen != "" {
			fuShenParts := strings.Split(lineInfo.FuShen, ":")
			if len(fuShenParts) == 2 {
				fuShenGanzhi := fuShenParts[1]
				fuShenWuXing := GetWuXingFromGanZhi(fuShenGanzhi)
				feiShenWuXing := lineWuXing

				lineDetail += fmt.Sprintf("\n  伏神→%s", lineInfo.FuShen)

				// Month/Day relationship with Fu Shen
				if IsKe(monthWuXing, fuShenWuXing) {
					lineDetail += " 月建克"
				}
				if IsKe(dayWuXing, fuShenWuXing) {
					lineDetail += " 日辰克"
				}

				// Fei Shen - Fu Shen relationship
				feiToFuRelation := GetRelation(feiShenWuXing, fuShenWuXing)
				if feiToFuRelation == "Sheng" {
					lineDetail += " 飞生伏"
				} else if feiToFuRelation == "Ke" {
					lineDetail += " 飞克伏"
				}

				fuToFeiRelation := GetRelation(fuShenWuXing, feiShenWuXing)
				if fuToFeiRelation == "Sheng" {
					lineDetail += " 伏生飞"
				} else if fuToFeiRelation == "Ke" {
					lineDetail += " 伏克飞"
				}

				// Check if Fu Shen is in Xun Kong
				if CheckXunKong(fuShenGanzhi, ctx.DayXunKong) {
					lineDetail += " 旬空"
				}
			}
		}

		result.Details = append(result.Details, lineDetail)
	}

	result.Details = append(result.Details, "")

	// Advanced Phase 3: Yuan Shen / Ji Shen Interactions
	// Check other moving lines and their relationships with changed lines
	yongShenWuXing := GetWuXingFromGanZhi(result.YongShenYao.Ganzhi)
	for i, changed := range ctx.Changed {
		if changed {
			// This is a moving line
			otherInfo := guaInfo[i]
			otherWuXing := GetWuXingFromGanZhi(otherInfo.Ganzhi)

			// Build detail string for this moving line
			interactionDetail := fmt.Sprintf("动爻互动: %s %s (%s, %s)", otherInfo.Position, otherInfo.Ganzhi, otherInfo.LiuQin, otherWuXing)

			// Check relationship with its changed line
			if len(bianGuaInfoAll) > i {
				bianOtherInfo := bianGuaInfoAll[i]
				bianOtherWuXing := GetWuXingFromGanZhi(bianOtherInfo.Ganzhi)

				// Check moving -> changed relationship (动爻对变爻的关系)
				movingToChangedRelation := GetRelation(otherWuXing, bianOtherWuXing)
				// Check changed -> moving relationship (变爻对动爻的关系)
				changedToMovingRelation := GetRelation(bianOtherWuXing, otherWuXing)

				// Build transformation description
				interactionDetail += fmt.Sprintf(" 化 %s (%s, %s)", bianOtherInfo.Ganzhi, bianOtherInfo.LiuQin, bianOtherWuXing)

				// Check for Jin Shen / Tui Shen
				jinTui := CheckJinTui(otherInfo.Ganzhi, bianOtherInfo.Ganzhi)
				if jinTui != "" {
					jinTuiCn := "进神"
					if jinTui == "Tui Shen" {
						jinTuiCn = "退神"
					}
					interactionDetail += fmt.Sprintf(", %s", jinTuiCn)
				}

				// Add Sheng/Ke relationship with traditional term and clear direction
				if changedToMovingRelation == "Sheng" {
					// Changed generates moving: 回头生(变生动)
					interactionDetail += ", 回头生(变生动)"
				} else if changedToMovingRelation == "Ke" {
					// Changed controls moving: 回头克(变克动)
					interactionDetail += ", 回头克(变克动)"
				} else if movingToChangedRelation == "Sheng" {
					// Moving generates changed: 泄气(动生变)
					interactionDetail += ", 泄气(动生变)"
				} else if movingToChangedRelation == "Ke" {
					// Moving controls changed: 克变(动克变)
					interactionDetail += ", 克变(动克变)"
				}
			}

			// Add interaction with Use God (skip if this IS the Use God line)
			if i != result.YongShenIndex {
				relation := GetRelation(otherWuXing, yongShenWuXing)

				if relation == "Sheng" {
					interactionDetail += fmt.Sprintf(", 生 用神 (%s) -> 吉", yongShenWuXing)
					result.Details = append(result.Details, interactionDetail)
					// Boost strength judgment?
					if result.Strength == "弱" {
						result.Strength = "强 (原神生助)"
					}
				} else if relation == "Ke" {
					interactionDetail += fmt.Sprintf(", 克 用神 (%s) -> 凶", yongShenWuXing)
					result.Details = append(result.Details, interactionDetail)
					// Reduce strength judgment?
					if result.Strength == "强" {
						result.Strength = "弱 (忌神克制)"
					}
				} else {
					// No direct Sheng/Ke relationship with Use God, but still show the moving line
					interactionDetail += fmt.Sprintf(", 与用神 (%s) 无直接生克", yongShenWuXing)
					result.Details = append(result.Details, interactionDetail)
				}
			} else {
				// This IS the Use God line - just show its transformation
				result.Details = append(result.Details, interactionDetail)
			}
		}
	}

	type BranchSource struct {
		Zhi      string
		Source   string
		IsDay    bool
		IsMonth  bool
		IsDong   bool // Is a moving line in Ben Gua (Dong Yao)
		IsBian   bool // Is a transformed line (Bian Yao)
		IsAnDong bool // Is a static line activated by day/month clash
	}

	var allBranches []BranchSource
	allBranches = append(allBranches, BranchSource{Zhi: ctx.DayZhi, Source: "日建", IsDay: true})
	allBranches = append(allBranches, BranchSource{Zhi: ctx.MonthZhi, Source: "月建", IsMonth: true})

	monthWuXing = GetWuXing(ctx.MonthZhi)

	for i, info := range guaInfo {
		zhi := string([]rune(info.Ganzhi)[1])
		isMoving := false
		if len(ctx.Changed) > i && ctx.Changed[i] {
			isMoving = true
		}

		isAnDong := false
		if !isMoving {
			lineWuXing := GetWuXingFromGanZhi(info.Ganzhi)
			monthStr := GetMonthStrength(lineWuXing, monthWuXing)
			if IsStrong(monthStr) && IsChong(ctx.DayZhi, zhi) {
				isAnDong = true
			}
		}

		allBranches = append(allBranches, BranchSource{
			Zhi:      zhi,
			Source:   info.Position,
			IsDong:   isMoving,
			IsAnDong: isAnDong,
		})

		if isMoving && len(bianGuaInfoAll) > i {
			bianZhi := string([]rune(bianGuaInfoAll[i].Ganzhi)[1])
			allBranches = append(allBranches, BranchSource{
				Zhi:    bianZhi,
				Source: info.Position + "变",
				IsBian: true,
			})
		}
	}

	yongShenZhi := string([]rune(result.YongShenYao.Ganzhi)[1])
	yongShenWuXing = GetWuXingFromGanZhi(result.YongShenYao.Ganzhi)
	bureauInfluence := 0

	// Check San He
	triads := []struct {
		Branches []string
		Element  string
	}{
		{[]string{"申", "子", "辰"}, "水"},
		{[]string{"亥", "卯", "未"}, "木"},
		{[]string{"寅", "午", "戌"}, "火"},
		{[]string{"巳", "酉", "丑"}, "金"},
	}

	for _, triad := range triads {
		count := 0
		members := []BranchSource{}
		foundBranches := make(map[string]bool)

		for _, target := range triad.Branches {
			for _, b := range allBranches {
				if b.Zhi == target {
					members = append(members, b)
					foundBranches[target] = true
					count++
					break
				}
			}
		}

		if count == 3 {
			numDongAn := 0
			numDayMonth := 0
			numBian := 0
			var parts []string
			containsYongShen := false

			for _, m := range members {
				label := m.Source
				if m.IsAnDong {
					label += "/暗动"
					numDongAn++
				} else if m.IsDong {
					numDongAn++
				} else if m.IsDay || m.IsMonth {
					numDayMonth++
				} else if m.IsBian {
					numBian++
				}

				parts = append(parts, fmt.Sprintf("%s(%s)", m.Zhi, label))
				if m.Zhi == yongShenZhi {
					containsYongShen = true
				}
			}

			// Shi Ju Check:
			// 1. 3 Dong/AnDong
			// 2. 2 Dong/AnDong + 1 Day/Month
			// 3. 1 Dong/AnDong + 1 Bian + 1 Day/Month
			isShiJu := (numDongAn == 3) ||
				(numDongAn == 2 && numDayMonth >= 1) ||
				(numDongAn == 1 && numBian >= 1 && numDayMonth >= 1)

			if isShiJu {
				bureauDesc := fmt.Sprintf("三合%s实局: %s", triad.Element, strings.Join(parts, " "))
				result.Details = append(result.Details, bureauDesc)

				// Impact on Yong Shen (Major)
				if containsYongShen {
					bureauInfluence += 4
				} else if triad.Element == yongShenWuXing {
					bureauInfluence += 3
				} else if IsSheng(triad.Element, yongShenWuXing) {
					bureauInfluence += 2
				} else if IsKe(triad.Element, yongShenWuXing) {
					bureauInfluence -= 4
				}
			} else if numDongAn > 0 || numBian > 0 {
				refreshDesc := fmt.Sprintf("地支增强(%s之力): %s", triad.Element, strings.Join(parts, " "))
				result.Details = append(result.Details, refreshDesc)
				if triad.Element == yongShenWuXing || IsSheng(triad.Element, yongShenWuXing) {
					bureauInfluence += 1
				}
			}
		}
	}

	// Check San Hui
	meetings := []struct {
		Branches []string
		Element  string
	}{
		{[]string{"亥", "子", "丑"}, "水"},
		{[]string{"寅", "卯", "辰"}, "木"},
		{[]string{"巳", "午", "未"}, "火"},
		{[]string{"申", "酉", "戌"}, "金"},
	}

	for _, meeting := range meetings {
		count := 0
		members := []BranchSource{}
		for _, target := range meeting.Branches {
			for _, b := range allBranches {
				if b.Zhi == target {
					members = append(members, b)
					count++
					break
				}
			}
		}

		if count == 3 {
			numDongAn := 0
			numDayMonth := 0
			numBian := 0
			var parts []string
			containsYongShen := false

			for _, m := range members {
				label := m.Source
				if m.IsAnDong {
					label += "/暗动"
					numDongAn++
				} else if m.IsDong {
					numDongAn++
				} else if m.IsDay || m.IsMonth {
					numDayMonth++
				} else if m.IsBian {
					numBian++
				}

				parts = append(parts, fmt.Sprintf("%s(%s)", m.Zhi, label))
				if m.Zhi == yongShenZhi {
					containsYongShen = true
				}
			}

			isShiJu := (numDongAn == 3) ||
				(numDongAn == 2 && numDayMonth >= 1) ||
				(numDongAn == 1 && numBian >= 1 && numDayMonth >= 1)

			if isShiJu {
				bureauDesc := fmt.Sprintf("三会%s实局: %s", meeting.Element, strings.Join(parts, " "))
				result.Details = append(result.Details, bureauDesc)

				if containsYongShen {
					bureauInfluence += 5
				} else if meeting.Element == yongShenWuXing {
					bureauInfluence += 3
				} else if IsSheng(meeting.Element, yongShenWuXing) {
					bureauInfluence += 2
				} else if IsKe(meeting.Element, yongShenWuXing) {
					bureauInfluence -= 5
				}
			} else if numDongAn > 0 || numBian > 0 {
				refreshDesc := fmt.Sprintf("地支增强(%s之力): %s", meeting.Element, strings.Join(parts, " "))
				result.Details = append(result.Details, refreshDesc)
				if meeting.Element == yongShenWuXing || IsSheng(meeting.Element, yongShenWuXing) {
					bureauInfluence += 1
				}
			}
		}
	}

	// Adjust Strength if bureau influence is significant
	if bureauInfluence >= 3 {
		if result.Strength == "弱" {
			result.Strength = "强 (合局生助)"
		} else if result.Strength == "中平" {
			result.Strength = "强"
		}
	} else if bureauInfluence <= -3 {
		if result.Strength == "强" {
			result.Strength = "弱 (合局克制)"
		} else if result.Strength == "中平" {
			result.Strength = "弱"
		}
	}

	// Phase 4: Judgment & Timing
	judgment, judgmentDetails := JudgeJiXiong(result.Strength, ctx.Category, ctx.Gender)
	result.Judgment = judgment
	result.Details = append(result.Details, judgmentDetails...)

	// Timing
	timing := PredictTiming(result.YongShenYao, judgment, ctx.DayZhi)
	result.Details = append(result.Details, fmt.Sprintf("应期: %s", timing))

	return result, nil
}

// DetermineYongShen maps the category to the corresponding Liu Qin
func DetermineYongShen(category string, gender string) string {
	switch category {
	case CategoryCareer:
		return "官鬼"
	case CategoryWealth:
		return "妻财"
	case CategoryMarriage:
		// For men: Wife Wealth, For women: Officer Ghost
		if gender == "Female" {
			return "官鬼"
		}
		// Default to "妻财" for Male or unspecified
		return "妻财"
	case CategoryStudy:
		return "父母"
	case CategorySafety:
		return "子孙"
	case CategoryHealth:
		// If asking for self: Shi Yao (World Line) - this is special.
		// If asking for parents: Parents, etc.
		// Let's assume "Health" means "Self Health" -> World Line?
		// But "Use God" usually refers to a Liu Qin.
		// If asking for self, the Use God is the Shi Yao itself.
		// We might need to handle this special case.
		// For now, let's return "世爻" and handle it.
		return "世爻"
	default:
		return "世爻" // Default to Self
	}
}

// Helper to get Wu Xing from Earthly Branch
func GetWuXing(zhi string) string {
	// Simple map
	m := map[string]string{
		"子": "水", "亥": "水",
		"寅": "木", "卯": "木",
		"巳": "火", "午": "火",
		"申": "金", "酉": "金",
		"辰": "土", "戌": "土", "丑": "土", "未": "土",
	}
	if val, ok := m[zhi]; ok {
		return val
	}
	return ""
}

// GetRelation returns the relationship between A and B (Sheng, Ke, Tong, etc.)
func GetRelation(a, b string) string {
	if a == b {
		return "Tong"
	}
	if IsSheng(a, b) {
		return "Sheng"
	}
	if IsKe(a, b) {
		return "Ke"
	}
	if IsSheng(b, a) {
		return "Xie"
	} // B Sheng A -> A Xie B

	// Strictly directional: Only return Ke if A controls B.
	// If B controls A, it is not A Ke B.

	return "None"
}

func TranslateRelation(rel string) string {
	switch rel {
	case "Sheng":
		return "生"
	case "Ke":
		return "克"
	case "Tong":
		return "同"
	case "Xie":
		return "泄"
	case "Chong":
		return "冲"
	case "He":
		return "合"
	default:
		return rel
	}
}

// CheckJinTui checks for Forward (Jin Shen) or Backward (Tui Shen) progress
func CheckJinTui(benGanzhi, bianGanzhi string) string {
	// Jin Shen: Yin->Mao, Si->Wu, Shen->You, Hai->Zi (Same Wu Xing, Yang -> Yin?)
	// Actually:
	// Hai -> Zi (Water)
	// Yin -> Mao (Wood)
	// Si -> Wu (Fire)
	// Shen -> You (Metal)
	// Chou -> Chen? (Earth) - Earth is complex.
	// Standard Jin Shen:
	// Hai->Zi, Yin->Mao, Si->Wu, Shen->You, Chou->Chen, Chen->Wei, Wei->Xu, Xu->Chou?

	// Simplified check:
	pairs := map[string]string{
		"亥": "子", "寅": "卯", "巳": "午", "申": "酉",
		"丑": "辰", "辰": "未", "未": "戌", "戌": "丑",
	}

	benZhi := string([]rune(benGanzhi)[1]) // 2nd char
	bianZhi := string([]rune(bianGanzhi)[1])

	if pairs[benZhi] == bianZhi {
		return "Jin Shen"
	}

	// Tui Shen is reverse
	for k, v := range pairs {
		if v == benZhi && k == bianZhi {
			return "Tui Shen"
		}
	}

	return ""
}

func CheckXunKong(ganzhi, dayXunKong string) bool {
	zhi := string([]rune(ganzhi)[1])
	return strings.Contains(dayXunKong, zhi)
}

func IsChong(a, b string) bool {
	chongMap := map[string]string{
		"子": "午", "午": "子",
		"丑": "未", "未": "丑",
		"寅": "申", "申": "寅",
		"卯": "酉", "酉": "卯",
		"辰": "戌", "戌": "辰",
		"巳": "亥", "亥": "巳",
	}
	return chongMap[a] == b
}

// Helper to check if A produces B (Sheng)
func IsSheng(a, b string) bool {
	shengMap := map[string]string{
		"金": "水", "水": "木", "木": "火", "火": "土", "土": "金",
	}
	return shengMap[a] == b
}

// Helper to check if A controls B (Ke)
func IsKe(a, b string) bool {
	keMap := map[string]string{
		"金": "木", "木": "土", "土": "水", "水": "火", "火": "金",
	}
	return keMap[a] == b
}

// CheckLiuHe checks for Six Combinations (Liu He) and returns the description
func CheckLiuHe(a, b string) string {
	// Zi-Chou -> Earth
	if (a == "子" && b == "丑") || (a == "丑" && b == "子") {
		return "子丑合土"
	}
	// Yin-Hai -> Wood
	if (a == "寅" && b == "亥") || (a == "亥" && b == "寅") {
		return "寅亥合木"
	}
	// Mao-Xu -> Fire
	if (a == "卯" && b == "戌") || (a == "戌" && b == "卯") {
		return "卯戌合火"
	}
	// Chen-You -> Metal
	if (a == "辰" && b == "酉") || (a == "酉" && b == "辰") {
		return "辰酉合金"
	}
	// Si-Shen -> Water
	if (a == "巳" && b == "申") || (a == "申" && b == "巳") {
		return "巳申合水"
	}
	// Wu-Wei -> Earth (or Fire/Earth?) - Standard is Earth (sums to Earth) or Fire (Summer)
	// Usually "Wu Wei He Tu"
	if (a == "午" && b == "未") || (a == "未" && b == "午") {
		return "午未合土"
	}

	return ""
}

// GetYaoName converts index (0-5) and bit ("0"or"1") to Yao Name (e.g. "初九", "六二")
func GetYaoName(index int, val string) string {
	// Index 0 is Bottom (Line 1/初)
	// Index 5 is Top (Line 6/上)

	positions := []string{"初", "二", "三", "四", "五", "上"}
	yinYang := "六"
	if val == "1" {
		yinYang = "九"
	}

	if index < 0 || index > 5 {
		return ""
	}

	pos := positions[index]

	// Rules:
	// Line 1 (Index 0): Pos + YinYang (e.g. 初九)
	// Line 2-5: YinYang + Pos (e.g. 九二)
	// Line 6 (Index 5): Pos + YinYang (e.g. 上六)

	if index == 0 || index == 5 {
		return pos + yinYang
	}
	return yinYang + pos
}

// CheckLiuHai checks for Six Harms (Liu Hai)
func CheckLiuHai(a, b string) bool {
	liuHaiMap := map[string]string{
		"子": "未", "未": "子",
		"丑": "午", "午": "丑",
		"寅": "巳", "巳": "寅",
		"卯": "辰", "辰": "卯",
		"申": "亥", "亥": "申",
		"酉": "戌", "戌": "酉",
	}
	return liuHaiMap[a] == b
}

// CheckXing checks for Punishments (Xing)
func CheckXing(a, b string) string {
	// San Xing (Three Punishments) - Strict check requires 3 branches.
	// But here we are comparing 2 branches (Line vs Day/Month/OtherLine).
	// So we can only detect "Partial Xing" or "Xiang Xing" between 2 branches.

	// Yin-Si-Shen (Tiger-Snake-Monkey)
	if (a == "寅" && b == "巳") || (a == "巳" && b == "寅") {
		return "寅巳相刑"
	}
	if (a == "巳" && b == "申") || (a == "申" && b == "巳") {
		return "巳申相刑"
	}
	if (a == "申" && b == "寅") || (a == "寅" && b == "申") {
		return "寅申相刑"
	}

	// Chou-Xu-Wei (Ox-Dog-Sheep)
	if (a == "丑" && b == "戌") || (a == "戌" && b == "丑") {
		return "丑戌相刑"
	}
	if (a == "戌" && b == "未") || (a == "未" && b == "戌") {
		return "未戌相刑"
	}
	if (a == "未" && b == "丑") || (a == "丑" && b == "未") {
		return "丑未相刑"
	}

	// Zi-Mao (Rat-Rabbit) - Rude Punishment
	if (a == "子" && b == "卯") || (a == "卯" && b == "子") {
		return "子卯相刑"
	}

	// Zi Xing (Self Punishment)
	if a == b {
		if a == "辰" {
			return "辰辰自刑"
		}
		if a == "午" {
			return "午午自刑"
		}
		if a == "酉" {
			return "酉酉自刑"
		}
		if a == "亥" {
			return "亥亥自刑"
		}
	}

	return ""
}

// CheckSanHe checks if a list of branches forms a San He (Three Harmony) combination
// Returns the resulting Element if formed, empty string otherwise.
func CheckSanHe(branches []string) string {
	// Need to check if the list contains all 3 required branches.
	// Shen-Zi-Chen -> Water
	// Hai-Mao-Wei -> Wood
	// Yin-Wu-Xu -> Fire
	// Si-You-Chou -> Metal

	has := func(target string) bool {
		for _, b := range branches {
			if b == target {
				return true
			}
		}
		return false
	}

	if has("申") && has("子") && has("辰") {
		return "水"
	}
	if has("亥") && has("卯") && has("未") {
		return "木"
	}
	if has("寅") && has("午") && has("戌") {
		return "火"
	}
	if has("巳") && has("酉") && has("丑") {
		return "金"
	}

	return ""
}

// CheckSanHui checks if a list of branches forms a San Hui (Three Meeting) combination
func CheckSanHui(branches []string) string {
	// Hai-Zi-Chou -> Water (North)
	// Yin-Mao-Chen -> Wood (East)
	// Si-Wu-Wei -> Fire (South)
	// Shen-You-Xu -> Metal (West)

	has := func(target string) bool {
		for _, b := range branches {
			if b == target {
				return true
			}
		}
		return false
	}

	if has("亥") && has("子") && has("丑") {
		return "水"
	}
	if has("寅") && has("卯") && has("辰") {
		return "木"
	}
	if has("巳") && has("午") && has("未") {
		return "火"
	}
	if has("申") && has("酉") && has("戌") {
		return "金"
	}

	return ""
}

// Phase 2: Strength Analysis

// Strength Levels
const (
	LevelWang  = "旺" // Prosperous
	LevelXiang = "相" // Strong
	LevelXiu   = "休" // Resting
	LevelQiu   = "囚" // Trapped
	LevelSi    = "死" // Dead
)

// CalculateStrength determines the strength of a Yao
func CalculateStrength(yaoInfo GuaInfo, bianYaoInfo *GuaInfo, isMoving bool, monthZhi, dayZhi, dayXunKong string) (string, []string) {
	details := []string{}
	yaoWuXing := GetWuXingFromGanZhi(yaoInfo.Ganzhi)
	yaoZhi := string([]rune(yaoInfo.Ganzhi)[1])

	monthWuXing := GetWuXing(monthZhi)
	dayWuXing := GetWuXing(dayZhi)

	score := 0

	// 1. Month Influence (Greatest)
	monthStrength := GetMonthStrength(yaoWuXing, monthWuXing)
	details = append(details, fmt.Sprintf("月建 (%s): %s", monthWuXing, monthStrength))
	if IsStrong(monthStrength) {
		score += 2
	}

	// 2. Day Influence (Second Greatest)
	dayStrength := GetDayStrength(yaoWuXing, dayWuXing)
	details = append(details, fmt.Sprintf("日辰 (%s): %s", dayWuXing, dayStrength))
	if IsStrong(dayStrength) {
		score += 2
	}

	// 3. Moving Line / Changed Line Influence
	bianStrength := ""
	if isMoving && bianYaoInfo != nil {
		bianWuXing := GetWuXingFromGanZhi(bianYaoInfo.Ganzhi)
		relation := GetRelation(bianWuXing, yaoWuXing)

		// Check for Jin Shen / Tui Shen
		jinTui := CheckJinTui(yaoInfo.Ganzhi, bianYaoInfo.Ganzhi)
		if jinTui != "" {
			jinTuiCn := "进神"
			if jinTui == "Tui Shen" {
				jinTuiCn = "退神"
			}
			details = append(details, fmt.Sprintf("变爻 (%s): %s (%s)", bianWuXing, TranslateRelation(relation), jinTuiCn))
			bianStrength = jinTui
		} else {
			details = append(details, fmt.Sprintf("变爻 (%s): %s", bianWuXing, TranslateRelation(relation)))
			if relation == "Sheng" {
				bianStrength = "Hui Tou Sheng"
			} else if relation == "Ke" {
				bianStrength = "Hui Tou Ke"
			} else if relation == "Xie" {
				bianStrength = "Xie Qi"
			}
		}
	}

	if bianStrength == "Hui Tou Sheng" || bianStrength == "Jin Shen" {
		score += 3
	} else if bianStrength == "Hui Tou Ke" || bianStrength == "Tui Shen" {
		score -= 5
	} else if bianStrength == "Xie Qi" {
		score -= 2
	}

	// 4. Advanced Interactions (Chong, He, Hai, Xing) with Month/Day
	if IsChong(monthZhi, yaoZhi) {
		details = append(details, "月破 (月冲)")
		score -= 4
	}
	if he := CheckLiuHe(monthZhi, yaoZhi); he != "" {
		details = append(details, fmt.Sprintf("月合 (%s)", he))
		score += 2
	}
	if he := CheckLiuHe(dayZhi, yaoZhi); he != "" {
		details = append(details, fmt.Sprintf("日合 (%s)", he))
		score += 2
	}
	if CheckLiuHai(dayZhi, yaoZhi) {
		details = append(details, "日害 (六害)")
		score -= 1
	}
	if xing := CheckXing(dayZhi, yaoZhi); xing != "" {
		details = append(details, fmt.Sprintf("日刑 (%s)", xing))
		score -= 1
	}

	// 5. Xun Kong / Ri Po / An Dong
	special := []string{}
	if CheckXunKong(yaoInfo.Ganzhi, dayXunKong) {
		special = append(special, "旬空")
		score -= 1
	}

	if IsChong(dayZhi, yaoZhi) {
		if !isMoving {
			if IsStrong(monthStrength) {
				special = append(special, "暗动")
				score += 1
			} else {
				special = append(special, "日破")
				score -= 3
			}
		} else {
			special = append(special, "日冲")
			score -= 1
		}
	}

	if len(special) > 0 {
		details = append(details, fmt.Sprintf("特殊状态: %s", strings.Join(special, ", ")))
	}

	// Final Conclusion
	overall := "弱"
	if score > 0 {
		overall = "强"
	} else if score == 0 {
		overall = "中平"
	}

	return overall, details
}

func GetMonthStrength(yao, month string) string {
	if yao == month {
		return LevelWang
	}
	if IsSheng(month, yao) {
		return LevelXiang
	}
	if IsSheng(yao, month) {
		return LevelXiu
	}
	if IsKe(yao, month) {
		return LevelQiu
	}
	if IsKe(month, yao) {
		return LevelSi
	}
	return LevelSi // Fallback
}

func GetDayStrength(yao, day string) string {
	// Day influence is similar to Month but technically Day "Sheng" is also strong support.
	// For simplicity, reuse the logic but note that Day "Ke" is not "Si" (Dead) but "Ke" (Controlled).
	// But the 5 levels apply mostly to Month.
	// For Day, we usually say: Lin Ri (Same), De Sheng (Born by Day), Bei Ke (Controlled), etc.
	// Let's map to the same levels for consistency.
	return GetMonthStrength(yao, day)
}

// Phase 3: Judgment & Timing

// JudgeJiXiong determines if the outcome is Auspicious or Inauspicious
func JudgeJiXiong(yongShenStrength string, category string, gender string) (string, []string) {
	judgment := "凶"
	isStrong := yongShenStrength == "强" || strings.Contains(yongShenStrength, "强")

	if isStrong {
		judgment = "吉"
	} else if yongShenStrength == "中平" {
		judgment = "平"
	}

	details := []string{fmt.Sprintf("吉凶判断: %s (基于用神旺衰: %s)", judgment, yongShenStrength)}

	// Gender & Category Specific Refinements
	// if category == CategoryMarriage {
	// 	if gender == "Female" {
	// 		if isStrong {
	// 			details = append(details, "女性测婚: 用神(官鬼)旺相，主夫星得力，缘分稳固。")
	// 		} else {
	// 			details = append(details, "女性测婚: 用神(官鬼)衰弱，需提防感情冷淡或阻碍。")
	// 		}
	// 	} else {
	// 		if isStrong {
	// 			details = append(details, "男性测婚: 用神(妻财)旺相，主妻贤家富，感情和谐。")
	// 		} else {
	// 			details = append(details, "男性测婚: 用神(妻财)衰弱，可能暗示求财或感情不顺。")
	// 		}
	// 	}
	// }

	return judgment, details
}

// PredictTiming estimates the time of the event
func PredictTiming(yongShenGuaInfo GuaInfo, judgment string, dayZhi string) string {
	// Simple Logic V1:
	// If Ji -> When Use God is strong (Wang/Xiang) or Combined (He).
	// If Xiong -> When Use God is weak or Clashed (Chong).

	// For now, return a placeholder or simple "Month/Day" based on Use God's Wu Xing.
	yongShenWuXing := GetWuXingFromGanZhi(yongShenGuaInfo.Ganzhi)

	return fmt.Sprintf("事件可能应验于 %s 日/月", yongShenWuXing)
}

func IsStrong(level string) bool {
	return level == LevelWang || level == LevelXiang
}

// Phase 4: Summary & Recommendations

// GenerateReport creates a formatted string of the analysis
func GenerateReport(result AnalysisResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== 六爻解卦报告 ===\n"))
	sb.WriteString(fmt.Sprintf("用神: %s (爻位: %s)\n", result.YongShen, result.YongShenYao.Position))
	sb.WriteString(fmt.Sprintf("总体旺衰: %s\n", result.Strength))
	sb.WriteString(fmt.Sprintf("吉凶: %s\n", result.Judgment))
	sb.WriteString(fmt.Sprintf("应期预测: %s\n", result.Details[len(result.Details)-1])) // Last detail is usually timing or judgment

	sb.WriteString("\n--- 分析详情 ---\n")
	for _, detail := range result.Details {
		sb.WriteString(fmt.Sprintf("- %s\n", detail))
	}

	// sb.WriteString("\n--- 建议 ---\n")
	// if result.Judgment == "吉" {
	// 	sb.WriteString("卦象吉利，可以积极行动，充满信心。\n")
	// } else {
	// 	sb.WriteString("卦象不佳，建议谨慎行事，静待时机。\n")
	// }

	return sb.String()
}
