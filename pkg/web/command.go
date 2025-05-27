package web

import (
	"github.com/kataras/iris/v12"
)

func GetCommandPost(ictx iris.Context, cmd any) error {
	err := ictx.ReadBody(cmd)
	if err != nil {
		return err
	}
	return nil
}

func GetCommand(ictx iris.Context, cmd any) error {
	err := ictx.ReadBody(cmd)
	if err != nil {
		return err
	}
	return nil
}
