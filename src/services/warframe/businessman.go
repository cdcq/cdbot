package warframe

import (
	"fmt"
	"time"
)

func BusinessmanResponse() string {
	ret := "虚空商人："
	nowTime := time.Now().Unix() / 60
	offset := int64(28021)
	nowTime = (nowTime + offset) % 20160
	if nowTime < 2880 {
		nowTime = 2880 - nowTime
		ret += "抵达\n"
	} else {
		nowTime = 20160 - nowTime
		ret += "离开\n"
	}
	day := nowTime / 1440
	hour := nowTime / 60 % 24
	minute := nowTime % 60
	ret += fmt.Sprintf("剩余 %02d 天 %02d 时 %02d 分", day, hour, minute)
	return ret
}
