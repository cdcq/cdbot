package api

import (
	"cdbot/src/global"
	"cdbot/src/models/requests"
	"cdbot/src/models/responses"
	"cdbot/src/services/warframe"
	"cdbot/src/services/xidian"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GenericMessageHandler(ctx *gin.Context) {
	var req requests.Message
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Abort()
		return
	}
	if req.MessageType != "group" {
		ctx.Abort()
		return
	}
	for _, group := range global.CONFIG.WFGroups {
		if group == req.GroupId {
			WFHandler(ctx, req)
			return
		}
	}
	for _, group := range global.CONFIG.XDGroups {
		if group == req.GroupId {
			XDHandler(ctx, req)
			return
		}
	}
}

func WFHandler(ctx *gin.Context, data requests.Message) {
	var res string
	message := data.RawMessage
	if message == "循环" {
		res = warframe.CycleResponse()
	} else if message == "黄历" {
		res = warframe.CalenderResponse()
	} else if message == "信条" {
		res = warframe.TenetResponse()
	} else if message == "奸商" {
		res = warframe.BusinessmanResponse()
	} else if len(message) > 13 && strings.HasPrefix(message, "信条更新\n") {
		err := warframe.TenetUpdate(message[13:])
		if err == nil {
			res = "1"
		} else {
			res = fmt.Sprintf("error: " + err.Error())
		}
	} else if len(message) > 3 && strings.HasPrefix(message, "wm ") {
		res = warframe.WMResponse(message[3:])
	} else {
		ctx.Abort()
		return
	}

	ok := responses.Result("group", 0, data.GroupId, res, false, ctx)

	if !ok {
		global.LOGGER.Warn("Send Message failed.")
	}
}

func XDHandler(ctx *gin.Context, data requests.Message) {
	var res string
	message := data.RawMessage
	if message == "黄历" {
		res = xidian.CalenderResponse()
	}
	responses.Result(data.MessageType, data.UserId, data.GroupId, res, false, ctx)
}
