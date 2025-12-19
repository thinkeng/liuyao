package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type YaoData struct {
	Index        int
	Name         string
	YaoCi        string
	XiangCi      string
	YaoDongHanYi string
}

type GuaData struct {
	Name        string
	BinaryCode  string
	GuaCi       string
	DaXiang     string
	CoreMeaning string
	Yaos        []YaoData
}

func main123() {
	content, err := os.ReadFile("卦辞.md")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	reGuaTitle := regexp.MustCompile(`####\s*\*?\s*([^：]+)：\s*([^（\s\*]+)\s*(\S*)\s*（([^）]+)`)
	reYaoBlockStart := regexp.MustCompile(`(\d+)\.\s*\*\*([^爻]+爻动)\s*（变([^）]+)`)
	reYaoDetail := regexp.MustCompile(`[\+\*->]\s*\*\*([^\*]+)\*\*：\s*(.+)`)

	allGua := make(map[string]*GuaData)
	var currentGua *GuaData
	var currentYao *YaoData

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := reGuaTitle.FindStringSubmatch(line); len(matches) > 4 {
			guaName := strings.Trim(strings.TrimSpace(matches[2]), "*")
			currentGua = &GuaData{
				Name:        guaName,
				CoreMeaning: strings.TrimSpace(matches[4]),
				Yaos:        make([]YaoData, 0),
			}
			allGua[guaName] = currentGua
			continue
		}
		if currentGua == nil {
			continue
		}

		if matches := reYaoBlockStart.FindStringSubmatch(line); len(matches) > 3 {
			idx, _ := strconv.Atoi(matches[1])
			currentYao = &YaoData{Index: idx, Name: strings.TrimSpace(matches[2])}
			continue
		}

		if matches := reYaoDetail.FindStringSubmatch(line); len(matches) > 2 {
			key := strings.TrimSpace(matches[1])
			val := strings.TrimSpace(matches[2])
			if key == "卦辞" {
				currentGua.GuaCi = val
			} else if currentYao != nil {
				switch key {
				case "本爻辞":
					currentYao.YaoCi = val
				case "爻动含义":
					currentYao.YaoDongHanYi = val
					currentGua.Yaos = append(currentGua.Yaos, *currentYao)
					currentYao = nil
				}
			}
		}
	}

	// 8宫 64卦 二进制 Key
	guaToBinary := map[string]string{
		"乾为天": "111111", "天风姤": "011111", "天山遁": "001111", "天地否": "000111", "风地观": "000011", "山地剥": "000001", "火地晋": "000101", "火天大有": "111101",
		"兑为泽": "110110", "泽水困": "010110", "泽地萃": "000110", "泽山咸": "001110", "水山蹇": "001010", "地山谦": "001000", "雷山小过": "001100", "雷泽归妹": "110100",
		"离为火": "101101", "火山旅": "001101", "火风鼎": "011101", "火水未济": "010101", "山水蒙": "010001", "风水涣": "010011", "天水讼": "010111", "天火同人": "101111",
		"震为雷": "100100", "雷地豫": "000100", "雷水解": "010100", "雷风恒": "011100", "地风升": "011000", "水风井": "011010", "泽风大过": "011110", "泽雷随": "100110",
		"巽为风": "011011", "风天小畜": "111011", "风火家人": "101011", "风雷益": "100011", "天雷无妄": "100111", "火雷噬嗑": "100101", "山雷颐": "100001", "山风蛊": "011001",
		"坎为水": "010010", "水泽节": "110010", "水雷屯": "100010", "水火既济": "101010", "泽火革": "101110", "雷火丰": "101100", "地火明夷": "101000", "地水师": "010000",
		"艮为山": "001001", "山火贲": "101001", "山天大畜": "111001", "山泽损": "110001", "火泽睽": "110101", "天泽履": "110111", "风泽中孚": "110011", "风山渐": "001011",
		"坤为地": "000000", "地雷复": "100000", "地泽临": "110000", "地天泰": "111000", "雷天大壮": "111100", "泽天夬": "111110", "水天需": "111010", "水地比": "000010",
	}

	// 完整的 大象辞 映射
	daXiang := map[string]string{
		"乾为天":  "天行健，君子以自强不息。",
		"坤为地":  "地势坤，君子以厚德载物。",
		"水雷屯":  "云雷，屯；君子以经纶。",
		"山水蒙":  "山下出泉，蒙；君子以果行育德。",
		"水天需":  "云上于天，需；君子以饮食宴乐。",
		"天水讼":  "天与水违行，讼；君子以作事谋始。",
		"地水师":  "地中有水，师；君子以容民畜众。",
		"水地比":  "地上有水，比；君子以建万国，亲诸侯。",
		"风天小畜": "风行天上，小畜；君子以懿文德。",
		"天泽履":  "上天下泽，履；君子以辨上下，定民志。",
		"地天泰":  "天地交，泰。后以财成天地之道，辅相天地之宜，以左右民。",
		"天地否":  "天地不交，否；君子以俭德辟难，不可荣以禄。",
		"天火同人": "天与火，同人；君子以类族辨物。",
		"火天大有": "火在天上，大有；君子以遏恶扬善，顺天休命。",
		"地山谦":  "地中有山，谦；君子以裒多益寡，称物平施。",
		"雷地豫":  "雷出地奋，豫。先王以作乐崇德，殷荐之上帝，以配祖考。",
		"泽雷随":  "泽中有雷，随；君子以向晦入宴息。",
		"山风蛊":  "山下有风，蛊；君子以振民育德。",
		"地泽临":  "泽上有地，临；君子以教思无穷，容保民无疆。",
		"风地观":  "风行地上，观。先王以省方观民设教。",
		"火雷噬嗑": "雷电噬嗑；先王以明罚敕法。",
		"山火贲":  "山下有火，贲；君子以明庶政，无敢折狱。",
		"山地剥":  "山附于地，剥；上以厚下安宅。",
		"地雷复":  "雷在地中，复。先王以至日闭关，商旅不行，后不省方。",
		"天雷无妄": "天下雷行，物与无妄；先王以茂对时育万物。",
		"山天大畜": "天在山中，大畜；君子以多识前言往行，以畜其德。",
		"山雷颐":  "山下有雷，颐；君子以慎言语，节饮食。",
		"泽风大过": "泽灭木，大过；君子以独立不惧，遁世无闷。",
		"坎为水":  "水洊至，习坎。君子以常德行，习教事。",
		"离为火":  "明两作，离；大人以继明照于四方。",
		"泽山咸":  "山上有泽，咸；君子以虚受人。",
		"雷风恒":  "雷风，恒；君子以立不易方。",
		"天山遁":  "天下有山，遁；君子以远小人，不恶而严。",
		"雷天大壮": "雷在天上，大壮；君子以非礼勿履。",
		"火地晋":  "明出地上，晋；君子以自昭明德。",
		"地火明夷": "明入地中，明夷；君子以莅众用晦而明。",
		"风火家人": "风自火出，家人；君子以言有物而行有恒。",
		"火泽睽":  "上火下泽，睽；君子以同而异。",
		"水山蹇":  "山上有水，蹇；君子以反身修德。",
		"雷水解":  "雷雨作，解；君子以赦过宥罪。",
		"山泽损":  "山下有泽，损；君子以惩忿窒欲。",
		"风雷益":  "风雷，益；君子以见善则迁，有过则改。",
		"泽天夬":  "泽上于天，夬；君子以施禄及下，居德则忌。",
		"天风姤":  "天下有风，姤；后以施命诰四方。",
		"泽地萃":  "泽上于地，萃；君子以除戎器，戒不虞。",
		"地风升":  "地中生木，升；君子以顺德，积小以高大。",
		"泽水困":  "泽无水，困；君子以致命遂志。",
		"水风井":  "木上有水，井；君子以劳民劝相。",
		"泽火革":  "泽中有火，革；君子以治历明时。",
		"火风鼎":  "木上有火，鼎；君子以正位凝命。",
		"震为雷":  "洊雷，震；君子以恐惧修省。",
		"艮为山":  "兼山，艮；君子以思不出其位。",
		"风山渐":  "山上有木，渐；君子以居贤德善俗。",
		"雷泽归妹": "泽上有雷，归妹；君子以永终知敝。",
		"雷火丰":  "雷电皆至，丰；君子以折狱致刑。",
		"火山旅":  "山上有火，旅；君子以明慎用刑，而不留狱。",
		"巽为风":  "随风，巽；君子以申命行事。",
		"兑为泽":  "丽泽，兑；君子以朋友讲习。",
		"风水涣":  "风行水上，涣。先王以享于帝立庙。",
		"水泽节":  "泽上有水，节；君子以制数度，议德行。",
		"风泽中孚": "泽上有风，中孚；君子以议狱缓死。",
		"雷山小过": "山上有雷，小过；君子以行过乎恭，丧过乎哀，用过乎俭。",
		"水火既济": "水在火上，既济；君子以思患而豫防之。",
		"火水未济": "火在水上，未济；君子以慎辨物居方。",
	}

	// 小象辞映射 (关键卦)
	xiaoXiang := map[string][]string{
		"乾为天": {"潜龙勿用，阳在下也。", "见龙在田，德施普也。", "终日乾乾，反复道也。", "或跃在渊，进无咎也。", "飞龙在天，大人造也。", "亢龙有悔，盈不可久也。"},
		"坤为地": {"履霜坚冰，阴始凝也。", "直以方也。不习无不利，地道光也。", "含章可贞，以时发也。", "括囊无咎，慎不害也。", "黄裳元吉，文在中也。", "战龙于野，其道穷也。"},
		"兑为泽": {"和兑之吉，行未疑也。", "孚兑之吉，信志也。", "来兑之凶，位不当也。", "商兑未宁，志不别也。", "孚于剥，恐疑惩也。", "引兑，未光也。"},
		"离为火": {"履错之敬，以辟咎也。", "黄离元吉，得中道也。", "日昃之离，何可久也？", "突如其来如，无所容也。", "六五之吉，离王公也。", "王用出征，以正邦也。"},
		"震为雷": {"震来虩虩，恐致福也。", "震来厉，乘刚也。", "震苏苏，位不当也。", "震遂泥，未光也。", "震往来厉，危行也。", "震索索，未得中也。"},
		"巽为风": {"进退，志疑也。", "纷若之吉，得中也。", "频巽之吝，志穷也。", "田获三品，富也。", "九五之吉，位正中也。", "巽在床下，上穷也。"},
		"坎为水": {"习坎入坎，失道凶也。", "坎有险，求小得。未出中也。", "来之坎坎，终无功也。", "樽酒簋贰，刚柔际也。", "坎不盈，中未大也。", "上六失道，凶三岁也。"},
		"艮为山": {"艮其趾，未失正也。", "艮其腓，未退听也。", "艮其限，危熏心也。", "艮其身，止诸躬也。", "艮其辅，以中正也。", "敦艮之吉，以厚终也。"},
		"泰卦":  {"拔茅征吉，志在外也。", "包荒得尚于中行，以光大也。", "无往不复，天地际也。", "翩翩不富，皆失实也。", "不戒以孚，中心愿也。", "城复于隍，其命乱也。"},
	}

	outFile, _ := os.Create("data/guadata.go")
	defer outFile.Close()
	w := bufio.NewWriter(outFile)
	fmt.Fprintln(w, "package data\n\nimport \"fmt\"\n\ntype YaoData struct {\n\tName         string\n\tYaoCi        string\n\tXiangCi      string\n\tYaoDongHanYi string\n}\n\ntype GuaData struct {\n\tName        string\n\tBinaryCode  string\n\tGuaCi       string\n\tDaXiang     string\n\tCoreMeaning string\n\tYaos        []YaoData\n}\n\nvar GuaIndex = map[string]GuaData{")

	for _, name := range []string{
		"乾为天", "天风姤", "天山遁", "天地否", "风地观", "山地剥", "火地晋", "火天大有",
		"兑为泽", "泽水困", "泽地萃", "泽山咸", "水山蹇", "地山谦", "雷山小过", "雷泽归妹",
		"离为火", "火山旅", "火风鼎", "火水未济", "山水蒙", "风水涣", "天水讼", "天火同人",
		"震为雷", "雷地豫", "雷水解", "雷风恒", "地风升", "水风井", "泽风大过", "泽雷随",
		"巽为风", "风天小畜", "风火家人", "风雷益", "天雷无妄", "火雷噬嗑", "山雷颐", "山风蛊",
		"坎为水", "水泽节", "水雷屯", "水火既济", "泽火革", "雷火丰", "地火明夷", "地水师",
		"艮为山", "山火贲", "山天大畜", "山泽损", "火泽睽", "天泽履", "风泽中孚", "风山渐",
		"坤为地", "地雷复", "地泽临", "地天泰", "雷天大壮", "泽天夬", "水天需", "水地比",
	} {
		binary := guaToBinary[name]
		g, ok := allGua[name]
		if !ok {
			continue
		}
		fmt.Fprintf(w, "\t\"%s\": {\n", binary)
		fmt.Fprintf(w, "\t\tName: \"%s\",\n\t\tBinaryCode: \"%s\",\n\t\tGuaCi: \"%s\",\n\t\tDaXiang: \"%s\",\n\t\tCoreMeaning: \"%s\",\n", g.Name, binary, g.GuaCi, daXiang[g.Name], g.CoreMeaning)
		fmt.Fprintln(w, "\t\tYaos: []YaoData{")
		for i, y := range g.Yaos {
			xc := ""
			if list, exists := xiaoXiang[name]; exists && i < len(list) {
				xc = list[i]
			}
			fmt.Fprintf(w, "\t\t\t{Name: \"%s\", YaoCi: \"%s\", XiangCi: \"%s\", YaoDongHanYi: \"%s\"},\n", y.Name, y.YaoCi, xc, y.YaoDongHanYi)
		}
		fmt.Fprintln(w, "\t\t},\n\t},")
	}
	fmt.Fprintln(w, "}\n\nfunc GetGuaData(binary string) (GuaData, bool) {\n\tg, ok := GuaIndex[binary]\n\treturn g, ok\n}\n\nfunc (g GuaData) Print() {\n\tfmt.Printf(\"【%s】(%s)\\n\", g.Name, g.BinaryCode)\n\tfmt.Printf(\"卦辞：%s\\n\", g.GuaCi)\n\tfmt.Printf(\"大象：%s\\n\", g.DaXiang)\n\tfor _, y := range g.Yaos {\n\t\tfmt.Printf(\"  %s：%s\\n\", y.Name, y.YaoCi)\n\t\tif y.XiangCi != \"\" {\n\t\t\tfmt.Printf(\"    《象》曰：%s\\n\", y.XiangCi)\n\t\t}\n\t}\n}")
	w.Flush()
}
