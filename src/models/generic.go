package models

import (
	"time"
)

type WmItem struct {
	ID       uint   `gorm:"primarykey"`
	ItemId   string `json:"id"`
	ItemName string `json:"item_name"`
	UrlName  string `json:"url_name"`
}

type WfMisc struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name"`
	Content   string `gorm:"type:text" json:"content"`
}
