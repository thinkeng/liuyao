// Package pkg provides functionality for Liu Yao divination.
package pkg

// GetShiYing returns the Shi (Subject) and Ying (Object) line positions (1-6) based on the index in the palace.
// indexInPalace: 0=Gua Wei, 1-5=1st-5th change, 6=You Hun, 7=Gui Hun
func GetShiYing(indexInPalace int) (shi, ying int) {
	switch indexInPalace {
	case 0: // 六冲卦 (本宫卦)
		return 6, 3
	case 1: // 一世卦
		return 1, 4
	case 2: // 二世卦
		return 2, 5
	case 3: // 三世卦
		return 3, 6
	case 4: // 四世卦
		return 4, 1
	case 5: // 五世卦
		return 5, 2
	case 6: // 游魂卦
		return 4, 1
	case 7: // 归魂卦
		return 3, 6
	default:
		return 0, 0
	}
}
