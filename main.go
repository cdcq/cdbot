package main

import (
	"cdbot/helpers"
	"cdbot/helpers/error_handlers"
	"cdbot/warframe"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	logFile, err := os.OpenFile("./logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	http.HandleFunc("/", receive)
	_ = http.ListenAndServe("127.0.0.1:5701", nil)
}

func receive(rw http.ResponseWriter, r *http.Request) {
	defer error_handlers.CloseHttpRequest(r)
	var data map[string]interface{}
	postData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.AddLog("main.go: receive", "read request body", err)
		return
	}
	decoder := json.NewDecoder(strings.NewReader(string(postData)))
	decoder.UseNumber()
	err = decoder.Decode(&data)
	if err != nil {
		helpers.AddLog("main.go: receive", "decode json", err)
		return
	}
	if data["message_type"] == "group" {
		message := data["message"].(string)
		groupId, err := data["group_id"].(json.Number).Int64()
		if err != nil {
			helpers.AddLog("main.go: receive", "turn json number to int64", err)
			return
		}
		if message == "给爷笑一个" {
			url := "http://127.0.0.1:5700/send_group_msg" +
				"?group_id=" + strconv.FormatInt(groupId, 10) +
				"&message=[CQ:face,id=13]"
			_, _ = http.Get(url)
			return
		}
		config := helpers.LoadConfig()
		if helpers.FindInI64Array(config.WFGroups, groupId) != -1 {
			warframe.WFHandler(data)
			return
		}
	}
}
