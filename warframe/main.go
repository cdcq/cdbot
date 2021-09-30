package warframe

import (
	"bytes"
	"cdbot/helpers"
	"encoding/json"
	"net/http"
	"strings"
)

func WFHandler(data map[string]interface{}) {
	funcName := "market.go: WMHandler"

	ret := make(map[string]interface{})
	groupId, err := data["group_id"].(json.Number).Int64()
	if err != nil {
		helpers.AddLog(funcName, "turn json number to int64", err)
		return
	}

	res := ""
	message := data["message"].(string)
	if message == "循环" {
		data["message"] = message[3:]
		res = CycleResponse()
	}
	if len(message) > 3 && strings.HasPrefix(message, "wm ") {
		data["message"] = message[3:]
		res = WMResponse(data["message"].(string))
	} else {
		return
	}

	ret["group_id"] = groupId
	ret["message"] = res
	retJson, err := json.Marshal(ret)
	if err != nil {
		helpers.AddLog(funcName, "marshal json", err)
		return
	}
	url := "http://127.0.0.1:5700/send_group_msg"
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(retJson))
}
