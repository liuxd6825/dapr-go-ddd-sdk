package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type CardFile struct {
	xbase.BaseModel `bson:",inline"`
	AppId           string `json:"appId" gorm:"app_id" bson:"app_id"`                      //应用Id
	FunId           string `json:"funId" gorm:"fun_id" bson:"fun_id"`                      //功能领域
	TitleText       string `json:"titleText" gorm:"title_text" bson:"title_text"`          //标题
	SubtitleText    string `json:"subtitleText" gorm:"subtitle_text" bson:"subtitle_text"` //子标题
	Icon            string `json:"icon" gorm:"icon" bson:"icon"`                           //图标
	Width           int    `json:"width" gorm:"width" bson:"width"`
	Height          int    `json:"height" gorm:"height" bson:"height"`
	CardClass       string `json:"cardClass" gorm:"card_class" bson:"card_class"`
	OrderNum        int64  `json:"orderNum" gorm:"order_num"  bson:"order_num"`
	CardHtml        string `json:"cardHtml" gorm:"card_html"  bson:"card_html"`
	LinkHtml        string `json:"linkHtml" gorm:"link_html"  bson:"link_html"`
}

func NewCardFile() (*CardFile, error) {
	return &CardFile{}, nil
}
