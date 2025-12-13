package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// YaoText 结构体：表示某一动爻的所有信息
type YaoText struct {
	Index        int    // 1-6
	YaoName      string // 例如: 初九, 六二, 上六
	BenYaoCi     string // 本爻辞
	BianGuaName  string // 变卦卦名
	BianGuaCi    string // 变卦卦辞
	YaoDongHanYi string // 爻动含义的完整描述
}

// GuaText 结构体：表示一个完整的卦象及其核心信息
type GuaText struct {
	Name        string             // 卦名 (例如：坤为地)
	Alias       string             // 别名/宫位名 (例如：本宫卦/一世卦)
	Hexagram    string             // 卦符 (例如：䷁)
	GuaCi       string             // 卦辞
	CoreMeaning string             // 核心意象
	ShiYao      string             // 世爻
	YaoMap      map[string]YaoText // 以 YaoName (初六, 九二) 为键的爻索引
}

// PalaceIndex 最终存储所有卦的索引：以卦名为键
var (
	PalaceIndex = make(map[string]*GuaText)
	once        sync.Once
)

// --- 正则表达式定义 ---
var (
	// reGuaTitle 修正：
	// Group 1: Alias部分 (如 **一、本宫卦**)
	// Group 2: 卦名 (如 坤为地)
	// Group 3: 卦符 (如 ䷁)
	// Group 4: 核心意象 (如 双重柔顺)
	reGuaTitle = regexp.MustCompile(`####\s*\*?\s*([^：]+)：\s*([^（\s]+)\s*(\S*)\s*（([^）]+)`)

	// 匹配卦辞、核心意象、世爻: + **XXX**：[内容] 或 * **XXX**：[内容]
	reProperty = regexp.MustCompile(`[\+\*]\s*\*\*([^\*]+)\*\*：\s*(.+)`)

	// 匹配爻动块的开始: 1. **初九爻动（变天风姤 ䷫）** -> 捕获: 1, 初九, 天风姤
	reYaoBlockStart = regexp.MustCompile(`(\d+)\.\s*\*\*([^爻]+爻动)\s*（变([^）]+)`)

	// 匹配爻辞、变卦辞、含义: - **XXX**：[内容] 或 * **XXX**：[内容]
	reYaoDetail = regexp.MustCompile(`[->\*]\s*\*\*([^\*]+)\*\*：\s*(.+)`)
)

// InitGuaCiIndex 解析 Markdown 文本，并建立结构化索引
func InitGuaCiIndex(markdownText string) {
	once.Do(func() {
		buildPalaceIndex(markdownText)
	})
}

// EnsureIndexInitialized ensures the index is built if not already.
// This might be useful if we load the file content from a specific location.
// For now, we assume the caller calls InitGuaCiIndex with content.
func buildPalaceIndex(markdownText string) {
	lines := strings.Split(markdownText, "\n")

	var currentGua *GuaText
	var currentYao YaoText

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 1. 匹配卦名 (遇到新卦则重置状态)
		if matches := reGuaTitle.FindStringSubmatch(line); len(matches) > 4 {
			if currentGua != nil && currentYao.YaoName != "" {
				currentGua.YaoMap[currentYao.YaoName] = currentYao
			}

			// 初始化新卦
			// 修复逻辑：对捕获的别名和卦名进行清理，去除可能残留的 * 字符
			alias := strings.Trim(strings.TrimSpace(matches[1]), "*")   // 别名
			guaName := strings.Trim(strings.TrimSpace(matches[2]), "*") // 卦名
			hexagram := strings.TrimSpace(matches[3])                   // 卦符
			coreMeaning := strings.TrimSpace(matches[4])                // 核心意象

			currentGua = &GuaText{
				Name:        guaName,
				Alias:       alias,
				Hexagram:    hexagram,
				CoreMeaning: coreMeaning,
				YaoMap:      make(map[string]YaoText),
			}
			PalaceIndex[guaName] = currentGua
			currentYao = YaoText{} // 重置爻状态
			continue
		}

		if line == "" || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "###") {
			continue
		}

		if currentGua == nil {
			continue
		}

		// 2. 匹配卦象核心属性 (卦辞、世爻、核心意象)
		if matches := reProperty.FindStringSubmatch(line); len(matches) > 2 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])

			switch key {
			case "卦辞":
				currentGua.GuaCi = value
			case "世爻":
				currentGua.ShiYao = value
			case "核心意象":
				if currentGua.CoreMeaning == "" {
					currentGua.CoreMeaning = value
				}
			}
			continue
		}

		// 3. 匹配爻动块的开始 (开启新爻的解析)
		if matches := reYaoBlockStart.FindStringSubmatch(line); len(matches) > 3 {
			if currentYao.YaoName != "" {
				currentGua.YaoMap[currentYao.YaoName] = currentYao
			}

			index, _ := strconv.Atoi(matches[1])
			currentYao = YaoText{
				Index:       index,
				YaoName:     strings.TrimSpace(matches[2]),
				BianGuaName: strings.TrimSpace(matches[3]),
			}
			continue
		}

		// 4. 匹配爻的详细信息 (本爻辞, 变卦辞, 爻动含义)
		if currentYao.YaoName != "" {
			if matches := reYaoDetail.FindStringSubmatch(line); len(matches) > 2 {
				key := strings.TrimSpace(matches[1])
				value := strings.TrimSpace(matches[2])

				switch key {
				case "本爻辞":
					currentYao.BenYaoCi = value
				case "变卦辞":
					currentYao.BianGuaCi = value
				case "爻动含义":
					currentYao.YaoDongHanYi = value
					currentGua.YaoMap[currentYao.YaoName] = currentYao
					currentYao = YaoText{}
				}
			}
		}
	}

	// 处理循环结束后的最后一爻（如果有）
	if currentGua != nil && currentYao.YaoName != "" {
		currentGua.YaoMap[currentYao.YaoName] = currentYao
	}
}

// QueryGuaAndYaoCi 根据卦名和动爻名查询信息
func QueryGuaAndYaoCi(guaName string, yaoName string) (GuaText, YaoText, error) {
	gua, ok := PalaceIndex[guaName]
	if !ok {
		return GuaText{}, YaoText{}, fmt.Errorf("错误：未找到卦名【%s】", guaName)
	}

	// 如果没有动爻，只返回卦的信息
	if yaoName == "" {
		return *gua, YaoText{}, nil
	}

	// 格式化查询的爻名，确保是 "初六爻动" 这样的全称
	yaoKey := yaoName
	if !strings.HasSuffix(yaoKey, "爻动") {
		yaoKey += "爻动"
	}

	yao, okYao := gua.YaoMap[yaoKey]
	if !okYao {
		// 尝试模糊匹配，比如输入 "初六"，匹配 "初六爻动"
		// 其实上面的 HasSuffix 处理了一部分，但如果 key 本身就是 "初六"，则可能不匹配
		// 不过根据 regex `reYaoBlockStart`，我们存进去的 key 是 matches[2] ("初九爻动")
		return *gua, YaoText{}, fmt.Errorf("已找到【%s】卦，但未找到动爻【%s】", guaName, yaoName)
	}

	return *gua, yao, nil
}
