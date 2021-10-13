package warframe

import (
	"cdbot/src/global"
	"cdbot/src/models"
	"errors"
	"gorm.io/gorm"
)

func TenetResponse() string {
	var item models.WfMisc
	if errors.Is(global.DATABASE.Where("name = ?", "tenet").First(&item).Error,
		gorm.ErrRecordNotFound) {
		return ""
	}
	return item.Content
}

func TenetUpdate(content string) error {
	item := models.WfMisc{
		Name:    "tenet",
		Content: content,
	}
	err := global.DATABASE.Updates(&item).Error
	return err
}
