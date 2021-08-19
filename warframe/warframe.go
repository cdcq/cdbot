package warframe

import (
	"bytes"
	"cdbot/helpers"
	"cdbot/helpers/error_handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func WMHandler(data map[string]interface{}) {
	ret := make(map[string]interface{})
	groupId, err := data["group_id"].(json.Number).Int64()
	if err != nil {
		helpers.AddLog("warframe.go: WMHandler", "turn json number to int64", err)
		return
	}
	if groupId != 692599380 {
		return
	}

	ret["group_id"] = groupId
	ret["message"] = WMResponse(data)
	retJson, err := json.Marshal(ret)
	if err != nil {
		helpers.AddLog("warframe.go: WMHandler", "marshal json", err)
		return
	}
	url := "http://127.0.0.1:5700/send_group_msg"
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(retJson))
}

func WMResponse(data map[string]interface{}) string {
	name := data["message"].(string)
	if name == "name" {
		nickNames, err := ioutil.ReadFile("./warframe/nick_names.yaml")
		if err != nil {
			helpers.AddLog("warframe.go: WMResponse", "read nick_names.yaml", err)
			return "出错了 :("
		}
		return string(nickNames)
	}
	rows := findWMItem(name)
	if rows == nil || len(rows) == 0 {
		return "啥也没找到\n" +
			"请尝试输入部分名称来获取提示"
	}
	if len(rows) == 1 {
		return parseWMResponse(rows[0].UrlName)
	}
	if len(rows) < 10 {
		isASet := true
		for i := 1; i < len(rows); i++ {
			if !strings.HasPrefix(rows[i].Name, rows[0].Name) {
				isASet = false
				break
			}
		}
		if isASet {
			return parseWMResponse(rows[0].UrlName)
		}
	}
	if len(rows) > 1 {
		ret := "你找的可能是：\n"
		for i, j := range rows {
			if i >= 10 {
				break
			}
			ret += j.Name + "\n"
		}
		if ret[len(ret) - 1] == '\n' {
			ret = ret[:len(ret) - 1]
		}
		return ret
	}
	return "出问题了 :("
}

type wmData struct {
	OrderType string
	Platinum int
	Quantity int
	Reputation int
}
type wmDataSlice []wmData
func (a wmDataSlice) Len() int {
	return len(a)
}
func (a wmDataSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a wmDataSlice) Less(i, j int) bool {
	if a[i].Platinum != a[j].Platinum {
		return a[i].Platinum < a[j].Platinum
	} else if a[i].Reputation != a[j].Reputation {
		return a[i].Reputation > a[j].Reputation
	} else {
		return a[i].Quantity > a[j].Quantity
	}
}

func parseWMResponse(urlName string) string {
	data := requireWMData(urlName)
	sort.Sort(wmDataSlice(data))
	ret := "查找" + urlName + "的结果：\n" +
		"白鸡 数量 名声 \n"
	cnt := 0
	for _, j := range data {
		cnt++
		if cnt >= 15 {
			break
		}
		ret += fmt.Sprintf("%-5d", j.Platinum) +
			fmt.Sprintf("%-5d", j.Quantity) +
			fmt.Sprintf("%-5d", j.Reputation) + "\n"
	}
	ret += "白鸡 数量 名声 "
	return ret
}

func requireWMData(urlName string) []wmData {
	url := "https://warframe.market/items/" + urlName
	res, err := http.Get(url)
	if err != nil {
		helpers.AddLog("warframe.go: requireWMData", "http get", err)
		return nil
	}
	defer error_handlers.CloseHttpResponse(res)
	resHtmlBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		helpers.AddLog("warframe.go: requireWMData", "read http response body", err)
		return nil
	}

	resHtml := string(resHtmlBytes)
	scriptStart := "<script type=\"application/json\" id=\"application-state\">"
	if strings.Index(resHtml, scriptStart) == -1 {
		return nil
	}
	jsonStart := strings.Index(resHtml, scriptStart) + len(scriptStart)
	dataJson := resHtml[jsonStart:]
	if !strings.HasPrefix(dataJson, "{\"payload\": {\"orders\": ") {
		return nil
	}
	jsonEnd := strings.Index(dataJson, "</script>")
	if strings.Index(resHtml, "</script>") == -1 {
		return nil
	}
	dataJson = dataJson[:jsonEnd]
	var data map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(dataJson))
	decoder.UseNumber()
	err = decoder.Decode(&data)
	if err != nil {
		helpers.AddLog("warframe.go: requireWMData", "decode json", err)
		return nil
	}
	data = data["payload"].(map[string]interface{})
	rows := data["orders"].([]interface{})

	var ret []wmData
	for _, j := range rows {
		row := j.(map[string]interface{})
		if row["platform"] != "pc" {
			continue
		}
		if row["order_type"] == "sell" {
			platinum, err := row["platinum"].(json.Number).Int64()
			if err != nil {
				helpers.AddLog("warframe.go: requireWMData", "turn json number to int64", err)
				return nil
			}
			quantity, err := row["quantity"].(json.Number).Int64()
			if err != nil {
				helpers.AddLog("warframe.go: requireWMData", "turn json number to int64", err)
				return nil
			}
			reputation, err := row["user"].
				(map[string]interface{})["reputation"].(json.Number).Int64()
			retRow := wmData{
				OrderType: "sell",
				Platinum: int(platinum),
				Quantity: int(quantity),
				Reputation: int(reputation),
			}
			ret = append(ret, retRow)
		}
	}
	return ret
}

type wmItem struct {
	Id string
	Name string
	UrlName string
	NickName string
}

func findWMItem(name string) []wmItem {
	db, err := sql.Open("sqlite3", "./warframe/database.db")
	if err != nil {
		helpers.AddLog("warframe.go: findWMItem", "open database", err)
		return nil
	}
	prep, err := db.Prepare("SELECT * FROM WM_ITEMS WHERE LOWER(NAME) LIKE LOWER(?) OR LOWER(NICK_NAME) LIKE LOWER(?) ORDER BY NAME ASC")
	if err != nil {
		helpers.AddLog("warframe.go: findWMItem", "prepare query", err)
		return nil
	}
	rows, err := prep.Query("%" + name + "%", "%" + name + "%")
	if err != nil {
		helpers.AddLog("warframe.go: findWMItem", "execute query", err)
		return nil
	}
	defer error_handlers.CloseRows(rows)
	var ret []wmItem
	for rows.Next() {
		var row wmItem
		err = rows.Scan(&row.Id, &row.Name, &row.UrlName, &row.NickName)
		if err != nil {
			helpers.AddLog("warframe.go: findWMItem", "scan rows", err)
			return nil
		}
		ret = append(ret, row)
	}
	return ret
}
