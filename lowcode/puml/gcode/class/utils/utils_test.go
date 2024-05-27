package utils

import "testing"

func TestRemoveTags(t *testing.T) {
	text := RemoveTags("截止日期 @enum {IdCard:身份证,Drive:驾驶证,Health:医保,SocialSecurity:社保卡} @data{}")
	t.Log(text)
}
