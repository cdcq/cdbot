package warframe

import "strings"

func WFHandler(data map[string]interface{}) {
	message := data["message"].(string)
	if message == "循环" {
		data["message"] = message[3:]
		CycleHander(data)
		return
	}
	if len(message) > 3 && strings.HasPrefix(message, "wm ") {
		data["message"] = message[3:]
		WMHandler(data)
		return
	}
}
