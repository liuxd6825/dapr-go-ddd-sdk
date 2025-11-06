package xtest

import "time"

func InitTimeZone() {
	local, err := time.LoadLocation("asia/shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = local
}
