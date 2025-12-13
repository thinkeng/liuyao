package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/thinkeng/liuyao/pkg"
)

func main3() {
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

	// Note: We can't easily access the length of the private map in pkg,
	// but we trust it's done. Or we could export a Len() function.
	// For now, let's just print a success message.
	fmt.Println("✅ 易经八宫数据索引建立完成。")
	fmt.Println(strings.Repeat("=", 60))

	// 3. 启动交互式查询
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入要查询的卦名和动爻名（例如：坤为地 初六）。输入 '退出' 结束程序。")

	for {
		fmt.Print("\n查询> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "退出" || input == "exit" {
			fmt.Println("程序结束。")
			break
		}

		parts := strings.Fields(input)
		if len(parts) != 2 {
			fmt.Println("输入格式错误。请按 [卦名] [动爻名] 格式输入。")
			continue
		}

		guaName := parts[0]
		yaoName := parts[1]

		gua, yao, err := pkg.QueryGuaAndYaoCi(guaName, yaoName)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(strings.Repeat("-", 40))
			fmt.Printf("【卦名】: %s %s (%s)\n", gua.Name, gua.Hexagram, gua.Alias)
			fmt.Printf("【卦辞】: %s\n", gua.GuaCi)
			fmt.Printf("【动爻】: %s (变 %s)\n", yao.YaoName, yao.BianGuaName)
			fmt.Printf("【本爻辞】: %s\n", yao.BenYaoCi)
			fmt.Printf("【爻动含义】: %s\n", yao.YaoDongHanYi)
			fmt.Println(strings.Repeat("-", 40))
		}
	}
}
