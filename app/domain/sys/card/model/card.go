package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/enum"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Card struct {
	xbase.BaseModel `bson:",inline"`
	HomeId          string            `json:"homeId" gorm:"home_id" bson:"home_id"`
	GroupId         string            `json:"groupId" gorm:"group_id" bson:"group_id"`
	CardId          string            `json:"cardId" gorm:"card_id" bson:"card_id"`
	TitleText       string            `json:"titleText" gorm:"title_text" bson:"title_text"`          //标题
	SubtitleText    string            `json:"subtitleText" gorm:"subtitle_text" bson:"subtitle_text"` //子标题
	Icon            string            `json:"icon" gorm:"icon" bson:"icon"`                           //图标
	Url             string            `json:"url" gorm:"url" bson:"url"`                              //链接
	OpenedTarget    enum.OpenedTarget `json:"openedTarget" gorm:"opened_target" bson:"opened_target"` //打开方式
	Width           int               `json:"width" gorm:"width" bson:"width"`
	Height          int               `json:"height" gorm:"height" bson:"height"`
	CardLevel       enum.CardLevel    `json:"cardLevel" gorm:"card_level" bson:"card_level"` //卡片级别
	CardClass       string            `json:"cardClass" gorm:"card_class" bson:"card_class"`
	OrderNum        int64             `json:"orderNum" gorm:"order_num"  bson:"order_num"`
	CardMode        enum.CardMode     `json:"cardMode" gorm:"card_mode"  bson:"card_mode"` //卡片模式 卡片或链接
	CardHtml        string            `json:"cardHtml" gorm:"card_html"  bson:"card_html"`
	LinkHtml        string            `json:"linkHtml" gorm:"link_html"  bson:"link_html"`
}

func NewCard() (*Card, error) {
	return &Card{}, nil
}
