package responses

import (
	"bytes"
	"cdbot/src/global"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

type MessagePack struct {
	MessageType string `json:"message_type"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id"`
	Message     string `json:"message"`
	AutoEscape  bool   `json:"auto_escape"`
}

type OkResp struct{}

func Result(messageType string, userId int64, groupId int64, message string, autoEscape bool, c *gin.Context) bool {
	c.JSON(http.StatusOK, OkResp{})
	jsonStr, err := json.Marshal(MessagePack{
		MessageType: messageType,
		UserId:      userId,
		GroupId:     groupId,
		Message:     message,
		AutoEscape:  autoEscape,
	})
	if err != nil {
		global.LOGGER.Warn("Message Type is wrong!")
		return false
	}
	req, err := http.NewRequest("POST", global.CONFIG.CQHttpAddr, bytes.NewBuffer(jsonStr))
	if err != nil {
		global.LOGGER.Warn("Could init http requests.")
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		global.LOGGER.Warn("Could not access CQHttp's API.")
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.LOGGER.Warn("Close body failed: " + err.Error())
			return
		}
	}(resp.Body)
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	return true
}
