package warframe

import (
	"cdbot/src/global"
	"cdbot/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func WMResponse(name string) string {
	name = strings.Replace(name, " ", "", -1)
	name = ProcessSpokenName(name)
	if name == "name" {
		nickNames, err := ioutil.ReadFile("./warframe/nick_names.yaml")
		if err != nil {
			global.LOGGER.Warn(err.Error())
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
			if !strings.HasPrefix(rows[i].ItemName, rows[0].ItemName) {
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
			ret += j.ItemName + "\n"
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
		global.LOGGER.Warn(funcName + "http get" + err.Error())
		return nil
	}

	resHtmlBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		global.LOGGER.Warn(funcName + "read http response body" + err.Error())
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
		global.LOGGER.Warn(funcName + "decode json" + err.Error())
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
				global.LOGGER.Warn(funcName + "turn json number to int64" + err.Error())
				return nil
			}
			quantity, err := row["quantity"].(json.Number).Int64()
			if err != nil {
				global.LOGGER.Warn(funcName + "turn json number to int64" + err.Error())
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

func findWMItem(name string) []models.WmItem {
	var ret []models.WmItem
	err := global.DATABASE.Model(&models.WmItem{}).Where("item_name LIKE ?", name).Find(&ret).Error
	if err != nil {
		global.LOGGER.Warn("Not found WMItems.")
		return nil
	}
	return ret
}

func checkSetName(name string) bool {
	var item models.WmItem

	if !errors.Is(global.DATABASE.Where("item_name = ?", name+"一套").
		First(&item).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
