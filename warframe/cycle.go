package warframe

import (
	"bytes"
	"cdbot/helpers"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CycleHander(data map[string]interface{}) {
	funcName := "cycle.go: CycleHandler"
	ret := make(map[string]interface{})
	groupId, err := data["group_id"].(json.Number).Int64()
	if err != nil {
		helpers.AddLog(funcName, "turn json number to int64", err)
		return
	}
	if groupId != 692599380 && groupId != 681638989 {
		return
	}

	ret["group_id"] = groupId
	ret["message"] = CycleResponse()
	retJson, err := json.Marshal(ret)
	if err != nil {
		helpers.AddLog(funcName, "marshal json", err)
		return
	}
	url := "http://127.0.0.1:5700/send_group_msg"
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(retJson))
}

func CycleResponse() string {
	funcName := "cycle.go: CycleResponse"

	ret := ""
	loc, err := time.LoadLocation("EST")
	if err != nil {
		helpers.AddLog(funcName, "load location", err)
		return "出错了 >n<"
	}
	nowTime := time.Now().Unix()

	ret = ret + "地球："
	hour := time.Now().In(loc).Hour()
	hour = (hour + 1) % 24
	minute := time.Now().In(loc).Minute()
	status, hour, minute := getEarthTime(hour, minute)
	if status == 0 {
		ret = ret + "白天"
	} else {
		ret = ret + "黑夜"
	}
	ret = ret + "\n" +
		fmt.Sprintf("剩余 %02d 时 %02d 分\n", hour, minute)

	ret = ret + "金星："
	status, minute, second := getVenusTime(nowTime)
	if status == 0 {
		ret = ret + "寒冷"
	} else {
		ret = ret + "温暖"
	}
	ret = ret + "\n" +
		fmt.Sprintf("剩余 %02d 分 %02d 秒\n", minute, second)
	return ret
}

func getEarthTime(hour, minute int) (int, int, int) {
	status := 0
	if hour%8 >= 4 {
		status = 1
	}
	hour = 3 - hour%4
	minute = 60 - minute
	return status, hour, minute
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
