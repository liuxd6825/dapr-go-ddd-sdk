package tokenutils

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/auth"
	"testing"
)

func TestGetLoginUser(t *testing.T) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJEdXhtLUp3dC1Ub2tlbiIsImV4cCI6MTcwMzc1Mjk3NCwidXNlciI6eyJ0ZW5hbnROYW1lIjoidGVzdCIsIm5hbWUiOiJ0ZXN0IiwidGVuYW50SWQiOiJ0ZXN0IiwidGVuYW50QWNjb3VudCI6InRlc3QiLCJpZCI6IjE3MjY0Nzk0NDEyNTUwNTEyNjQiLCJ1c2VyVHlwZSI6IlRFTkFOVF9BRE1JTiIsImFjY291bnQiOiJ0ZXN0Iiwic3RhdHVzIjoiVVNFSU5HIn0sImNsaWVudF9pZCI6IjA5OGY2YmNkNDYyMWQzNzNjYWRlNGU4MzI2MjdiNGY2In0.s_kHa3pKt6XehbsL7E9PJqywM_pxbbq6V2zHyZCJmDk"
	ctx, err := auth.NewContext(context.Background(), token, "")
	user, err := auth.GetLoginUser(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("userId=%s, userName=%s", user.GetId(), user.GetName())
}
