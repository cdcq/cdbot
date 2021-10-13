package warframe

import (
	"fmt"
	"time"
)

func CycleResponse() string {
	// funcName := "cycle.go: CycleResponse"

	ret := ""
	nowTime := time.Now().Unix()

	/*
		loc, err := time.LoadLocation("EST")
		if err != nil {
			helpers.AddLog(funcName, "load location", err)
			return "出错了 >n<"
		}
		ret += "地球："
		hour := time.Now().In(loc).Hour()
		hour = (hour + 1) % 24
		minute := time.Now().In(loc).Minute()
		status, hour, minute := getEarthTime(hour, minute)
		if status == 0 {
			ret += "黑夜"
		} else {
			ret += "白天"
		}
		ret += "\n" +
			fmt.Sprintf("剩余 %02d 时 %02d 分\n", hour, minute)
	*/

	ret += "夜灵平野："
	status, minute, second := getCetusTime(nowTime)
	if status == 1 {
		ret += "白昼"
	} else {
		ret += "黑夜"
	}
	ret += "\n" +
		fmt.Sprintf("剩余 %02d 分 %02d 秒\n", minute, second)

	ret += "奥布山谷："
	status, minute, second = getVenusTime(nowTime)
	if status == 0 {
		ret += "寒冷"
	} else {
		ret += "温暖"
	}
	ret += "\n" +
		fmt.Sprintf("剩余 %02d 分 %02d 秒\n", minute, second)
	return ret
}

/*
func getEarthTime(hour, minute int) (int, int, int) {
	status := 0
	if hour%8 >= 4 {
		status = 1
	}
	hour = 3 - hour%4
	minute = 60 - minute
	return status, hour, minute
}
*/

func getCetusTime(nowTime int64) (int, int, int) {
	offset := int64(-1380)
	nowTime = (nowTime + offset) % 9000
	status := 0 // 0 is day, 1 is night.
	if nowTime <= 6000 {
		nowTime = int64(6000) - nowTime
	} else {
		nowTime = int64(3000) - (nowTime - 6000)
		status = 1
	}
	return status, int(nowTime / 60), int(nowTime % 60)
}

func getVenusTime(nowTime int64) (int, int, int) {
	offset := int64(1226)
	nowTime = (nowTime - offset) % 1600
	status := 0 // 0 is cold, 1 is warm.
	if nowTime <= 1200 {
		nowTime = int64(1200) - nowTime
	} else {
		nowTime = int64(400) - (nowTime - 1200)
		status = 1
	}
	return status, int(nowTime / 60), int(nowTime % 60)
}
