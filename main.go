package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/thinkeng/liuyao/pkg"
)

// // 定义爻信息
// type YaoInfo struct {
// 	position string // 爻位 (初、二、三、四、五、上)
// 	ben      string // 本卦阴阳 (阳爻"─", 阴爻"--")
// 	bian     string // 变卦阴阳
// 	dong     string // 动爻标记 ("○"动爻, "●"变爻, " "不动)
// 	naJia    string // 纳甲 (地支)
// 	shiYing  string // 世应标记 ("世"或"应")
// 	liuQin   string // 六亲
// 	liuShen  string // 六神
// }

// 生成随机投掷结果
func randomToss() string {
	// 生成三个随机数字（0或1）
	result := make([]byte, 3)
	for i := range result {
		if rand.Intn(2) == 1 {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}
	return string(result)
}

func main() {

	const filename = "卦辞.md"

	// 1. 读取文件内容
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ 文件读取失败：请确保文件【%s】存在于当前目录，并包含正确的八宫Markdown内容。\n错误信息: %v\n", filename, err)
		return
	}

	markdownContent := string(contentBytes)

	// 2. 建立结构体索引
	pkg.InitGuaCiIndex(markdownContent)

	//============

	tosses := make([]string, 6)
	for i := range tosses {

		tosses[i] = randomToss()
		fmt.Println(pkg.ParseToss(tosses[i]))
	}

	fmt.Println(tosses)
	// 生成卦
	gua := pkg.GenerateGua(tosses)

	// 打印卦象
	printGua(gua)

	baZi, dayKong := pkg.GetDayGanZhi(time.Now())
	fmt.Println("日期: ", baZi.GetYear()+" "+baZi.GetMonth()+" "+baZi.GetDay()+" "+baZi.GetTime())
	fmt.Println("旬空: ", dayKong)

	dayGan := baZi.GetDayGan()
	dayZhi := baZi.GetDayZhi()
	monthZhi := baZi.GetMonthZhi()

	hexagram := strings.Join(gua.BenGua, "")
	if result, err := pkg.GetGuaInfo(hexagram, dayGan); err == nil {
		guaName := pkg.DetermineGuaName(hexagram)
		fullGuaName := pkg.GetFullGuaName(hexagram)
		palaceIndex, _, _ := pkg.GetGuaPalace(guaName)
		palaceWuXing := pkg.GetPalaceWuXing(palaceIndex)
		fmt.Printf("本卦: %s (%s) 纳甲与六神配置 (日干:%s):\n", fullGuaName, hexagram, dayGan)

		// Display Shen Sha Config
		shenShaConfig := pkg.GetShenShaConfig(dayGan, dayZhi, monthZhi)
		fmt.Printf("神煞: %s\n", strings.Join(shenShaConfig, " "))

		fmt.Println("====================================")
		//fmt.Println("爻位\t干支\t六神\t六亲\t世应\t伏神\t爻类型")
		fmt.Println("爻位\t六神\t六亲\t干支\t伏神    \t世应\t爻类型")
		fmt.Println("------------------------------------")

		for i := len(result) - 1; i >= 0; i-- {
			info := result[i]
			// Override YaoType with specific coin toss result (Lao Yang/Lao Yin)
			specificType, specificName := pkg.ParseToss(tosses[i])
			info.YaoType = specificName + ":" + specificType

			fmt.Printf("%s\t%s\t%s\t%s\t%-8s\t%s\t%s\n", info.Position, info.LiuShen, info.LiuQin, info.Ganzhi, info.FuShen, info.ShiYing, info.YaoType)
		}
		fmt.Println("====================================")

		// Check if there are moving lines to display Bian Gua
		hasMoving := false
		for _, changed := range gua.Changed {
			if changed {
				hasMoving = true
				break
			}
		}

		if hasMoving {
			bianHexagram := strings.Join(gua.BianGua, "")
			// Use GetBianGuaInfo with Ben Gua's Palace Wu Xing
			if bianResult, err := pkg.GetBianGuaInfo(bianHexagram, dayGan, palaceWuXing); err == nil {
				bianGuaName := pkg.GetFullGuaName(bianHexagram)
				fmt.Printf("\n变卦: %s (%s) 纳甲与六神配置:\n", bianGuaName, bianHexagram)
				fmt.Println("====================================")
				//fmt.Println("爻位\t干支\t六神\t六亲\t世应\t爻类型")
				fmt.Println("爻位\t六神\t六亲\t干支\t伏神\t世应\t爻类型")
				fmt.Println("------------------------------------")

				for i := len(bianResult) - 1; i >= 0; i-- {
					info := bianResult[i]
					fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n", info.Position, info.LiuShen, info.LiuQin, info.Ganzhi, info.FuShen, info.ShiYing, info.YaoType)
				}
				fmt.Println("====================================")
			}
		}

	} else {
		fmt.Println("错误:", err)
	}

	// Hexagram Analysis
	fmt.Println("\n====================================")
	fmt.Println("开始解卦 (Hexagram Analysis)...")
	fmt.Println("====================================")

	// For now, hardcode a category or use a simple input simulation
	// In a real CLI, we would use flags or interactive input.
	// Let's default to "Wealth" for demonstration.
	//category := pkg.CategoryWealth
	category := pkg.CategoryMarriage
	fmt.Printf("设定求测事项: %s (默认)\n", category)

	// Need Bian Hexagram and Changed array
	// gua.BianGua is []string, need to join
	bianHexagram := strings.Join(gua.BianGua, "")

	analysisCtx := pkg.AnalysisContext{
		GuaHexagram:  hexagram,
		BianHexagram: bianHexagram,
		Changed:      gua.Changed,
		DayGan:       dayGan,
		DayZhi:       dayZhi,
		MonthZhi:     monthZhi,
		DayXunKong:   dayKong,
		Category:     category,
		Gender:       "Female", // Default for demo Female
		Date:         time.Now(),
	}

	analysisResult, err := pkg.Analyze(analysisCtx)
	if err != nil {
		fmt.Printf("解卦失败: %v\n", err)
	} else {
		report := pkg.GenerateReport(analysisResult)
		fmt.Println(report)

		// Display Text Info (Gua & Yao)
		fmt.Println("================动爻卦辞====================")
		guaText, _, _ := pkg.QueryGuaAndYaoCi(analysisResult.GuaName, "")

		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("【卦名】: %s %s (%s)\n", guaText.Name, guaText.Hexagram, guaText.Alias)
		fmt.Printf("【卦辞】: %s\n", guaText.GuaCi)

		for _, yao := range analysisResult.MovingYaos {
			fmt.Println(strings.Repeat("-", 20))
			fmt.Printf("【动爻】: %s (变 %s)\n", yao.YaoName, yao.BianGuaName)
			fmt.Printf("【本爻辞】: %s\n", yao.BenYaoCi)
			fmt.Printf("【爻动含义】: %s\n", yao.YaoDongHanYi)
		}
		fmt.Println(strings.Repeat("-", 40))
	}
}

// 打印卦象
func printGua(gua pkg.Gua) {
	fmt.Println("本卦 → 变卦:", strings.Join(gua.BenGua, ""), strings.Join(gua.BianGua, ""))
	//fmt.Println("爻象 (0=阴, 1=阳):", strings.Join(gua.yao, " "))

	fmt.Print("动爻: [")
	for i, dong := range gua.Changed {
		if dong {
			fmt.Printf("%s ", []string{"初", "二", "三", "四", "五", "上"}[i])
		}
	}
	fmt.Println("]")
	fmt.Println()
}
