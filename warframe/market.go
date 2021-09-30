package warframe

import (
	"cdbot/helpers"
	"cdbot/helpers/error_handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func WMResponse(name string) string {
	funcName := "market.go: WMResponse"

	name = strings.Replace(name, " ", "", -1)
	name = ProcessSpokenName(name)
	if name == "name" {
		nickNames, err := ioutil.ReadFile("./warframe/nick_names.yaml")
		if err != nil {
			helpers.AddLog(funcName, "read nick_names.yaml", err)
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
		if ret[len(ret)-1] == '\n' {
			ret = ret[:len(ret)-1]
		}
		return ret
	}
	return "出问题了 :("
}

func ProcessSpokenName(name string) string {
	name = strings.ToLower(name)

	var nickNames map[string][]string
	yamlFile, _ := ioutil.ReadFile("./warframe/nick_names.yaml")
	_ = yaml.Unmarshal(yamlFile, &nickNames)
	for key, value := range nickNames {
		for _, nickName := range value {
			name = strings.Replace(name, nickName+"甲", key+"prime", -1)
			name = strings.Replace(name, nickName, key+"prime", -1)
		}
		if name == strings.ToLower(key) {
			name = name + "prime一套"
			break
		}
	}

	name = strings.Replace(name, "总图", "图", -1)
	name = strings.Replace(name, "蓝图", "图", -1)
	name = strings.Replace(name, "图", "蓝图", -1)

	name = strings.Replace(name, "头部神经光元", "头", -1)
	name = strings.Replace(name, "头", "头部神经光元", -1)

	name = strings.Replace(name, "prime", "p", -1)
	name = strings.Replace(name, "pp", "p", -1)
	name = strings.Replace(name, "p", "prime", -1)

	if checkSetName(name) {
		name = name + "一套"
	}
	return name
}

type wmData struct {
	OrderType  string
	Platinum   int
	Quantity   int
	Reputation int
	Status     string
}
type wmDataSlice []wmData

func (a wmDataSlice) Len() int {
	return len(a)
}
func (a wmDataSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a wmDataSlice) Less(i, j int) bool {
	if a[i].Status == "ingame" && a[j].Status != "ingame" {
		return true
	}
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
		"白鸡 数量 名声 状态\n"
	cnt := 0
	for _, j := range data {
		if cnt < 5 && j.Status != "ingame" {
			continue
		}
		if cnt >= 5 && j.Status == "ingame" {
			continue
		}
		if cnt >= 10 {
			break
		}
		cnt++
		status := j.Status
		if j.Status == "ingame" {
			status = "在线"
		} else {
			status = "离线"
		}
		ret += fmt.Sprintf("%-5d", j.Platinum) +
			fmt.Sprintf("%-5d", j.Quantity) +
			fmt.Sprintf("%-5d", j.Reputation) +
			fmt.Sprintf("%-5s", status) + "\n"
	}
	return ret
}

func requireWMData(urlName string) []wmData {
	funcName := "market.go: requireWMData"

	url := "https://warframe.market/items/" + urlName
	res, err := http.Get(url)
	if err != nil {
		helpers.AddLog(funcName, "http get", err)
		return nil
	}
	defer error_handlers.CloseHttpResponse(res)
	resHtmlBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		helpers.AddLog(funcName, "read http response body", err)
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
		helpers.AddLog(funcName, "decode json", err)
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
				helpers.AddLog(funcName, "turn json number to int64", err)
				return nil
			}
			quantity, err := row["quantity"].(json.Number).Int64()
			if err != nil {
				helpers.AddLog(funcName, "turn json number to int64", err)
				return nil
			}
			reputation, err := row["user"].(map[string]interface{})["reputation"].(json.Number).Int64()
			status := row["user"].(map[string]interface{})["status"].(string)
			retRow := wmData{
				OrderType:  "sell",
				Platinum:   int(platinum),
				Quantity:   int(quantity),
				Reputation: int(reputation),
				Status:     status,
			}
			ret = append(ret, retRow)
		}
	}
	return ret
}

type wmItem struct {
	Id      string
	Name    string
	UrlName string
}

func findWMItem(name string) []wmItem {
	funcName := "market.go: findWMItem"

	db, err := sql.Open("sqlite3", "./warframe/database.db")
	if err != nil {
		helpers.AddLog(funcName, "open database", err)
		return nil
	}
	prep, err := db.Prepare("SELECT * FROM WM_ITEMS WHERE LOWER(NAME) LIKE LOWER(?)")
	if err != nil {
		helpers.AddLog(funcName, "prepare query", err)
		return nil
	}
	rows, err := prep.Query("%" + name + "%")
	if err != nil {
		helpers.AddLog(funcName, "execute query", err)
		return nil
	}
	defer error_handlers.CloseRows(rows)
	var ret []wmItem
	for rows.Next() {
		var row wmItem
		err = rows.Scan(&row.Id, &row.Name, &row.UrlName)
		if err != nil {
			helpers.AddLog(funcName, "scan rows", err)
			return nil
		}
		ret = append(ret, row)
	}
	return ret
}

func checkSetName(name string) bool {
	funcName := "market.go: checkSetName"

	db, err := sql.Open("sqlite3", "./warframe/database.db")
	if err != nil {
		helpers.AddLog(funcName, "open database", err)
		return false
	}
	prep, err := db.Prepare("SELECT * FROM WM_ITEMS WHERE LOWER(NAME)=LOWER(?)")
	if err != nil {
		helpers.AddLog(funcName, "prepare query", err)
		return false
	}
	rows, err := prep.Query(name + "一套")
	if err != nil {
		helpers.AddLog(funcName, "execute query", err)
		return false
	}
	defer error_handlers.CloseRows(rows)
	if rows.Next() {
		return true
	} else {
		return false
	}
}
