package pkg

var (
	// 爻位名称
	yaoPositions = []string{"初爻", "二爻", "三爻", "四爻", "五爻", "上爻"}

	// 爻类型映射
	yaoTypeMap = map[string]string{
		"111": "老阳:—○", "000": "老阴:⚋×",
		"110": "少阴:⚋", "001": "少阳:—", // 110=2 Heads(Yin), 001=1 Head(Yang)
		"1": "阳爻:—", "0": "阴爻:⚋",
	}
)
