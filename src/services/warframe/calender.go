package warframe

import (
	"cdbot/src/helpers"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"time"
)

type Words struct {
	Activity  []string `yaml:"Activity"`
	Direction []string `yaml:"Direction"`
	Place     []string `yaml:"Place"`
}

func CalenderResponse() string {
	funcName := "calender.go: CalenderResponse"

	unixDay := (time.Now().Unix() - 1633017600) / 86400
	r := rand.New(rand.NewSource(unixDay))

	words := Words{}
	configFile, err := ioutil.ReadFile("./warframe/calender_words.yaml")
	if err != nil {
		helpers.AddLog(funcName, "read file", err)
		return "Something Wrong. :("
	}
	err = yaml.Unmarshal(configFile, &words)
	if err != nil {
		helpers.AddLog(funcName, "unmarshal yaml", err)
		return "Something Wrong. :("
	}
	chosen := make([]string, 0)
	activities := words.Activity
	directions := words.Direction
	places := words.Place

	res := "今日运势\n"

	res += "宜：\n"
	n1 := 2 + r.Int()%2
	l1 := len(activities)
	for i := 1; i <= n1; i++ {
		chose := r.Intn(l1)
		for helpers.FindInStringArray(chosen, activities[chose]) != -1 {
			chose = r.Intn(l1)
		}
		chosen = append(chosen, activities[chose])
		res += activities[chose] + " "
	}
	res += "\n"

	res += "忌：\n"
	n2 := 2 + r.Int()%2
	for i := 1; i <= n2; i++ {
		chose := r.Intn(l1)
		for helpers.FindInStringArray(chosen, activities[chose]) != -1 {
			chose = r.Intn(l1)
		}
		chosen = append(chosen, activities[chose])
		res += activities[chose] + " "
	}
	res += "\n"

	res += "方位：\n" +
		fmt.Sprintf("面向[%s]方向输出伤害最高\n", directions[r.Intn(len(directions))]) +
		fmt.Sprintf("位于[%s]开核桃出金最多\n", places[r.Intn(len(places))]) +
		"\n仅供娱乐，不做参考"
	return res
}
