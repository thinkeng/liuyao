package pkg

import (
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// 获取日干支 (修正)
func GetDayGanZhi(date time.Time) (*calendar.EightChar, string) {
	// 阳历转八字
	year := date.Year()     // 年
	month := date.Month()   // 月 (time.Month 类型)
	day := date.Day()       // 日
	hour := date.Hour()     // 时
	minute := date.Minute() // 分
	second := date.Second() // 秒

	solar := calendar.NewSolar(year, int(month), day, hour, minute, second)
	lunar := solar.GetLunar()

	baZi := lunar.GetEightChar()
	//fmt.Println(baZi.GetYear() + " " + baZi.GetMonth() + " " + baZi.GetDay() + " " + baZi.GetTime())

	//return baZi.GetYear() + " " + baZi.GetMonth() + " " + baZi.GetDay() + " " + baZi.GetTime()
	return baZi, lunar.GetDayXunKong()
}
