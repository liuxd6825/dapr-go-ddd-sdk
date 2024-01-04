package restapp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContext(t *testing.T) {
	DefaultAuthToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJEdXhtLUp3dC1Ub2tlbiIsImV4cCI6MTcwMzc1Mjk3NCwidXNlciI6eyJ0ZW5hbnROYW1lIjoidGVzdCIsIm5hbWUiOiJ0ZXN0IiwidGVuYW50SWQiOiJ0ZXN0IiwidGVuYW50QWNjb3VudCI6InRlc3QiLCJpZCI6IjE3MjY0Nzk0NDEyNTUwNTEyNjQiLCJ1c2VyVHlwZSI6IlRFTkFOVF9BRE1JTiIsImFjY291bnQiOiJ0ZXN0Iiwic3RhdHVzIjoiVVNFSU5HIn0sImNsaWVudF9pZCI6IjA5OGY2YmNkNDYyMWQzNzNjYWRlNGU4MzI2MjdiNGY2In0.s_kHa3pKt6XehbsL7E9PJqywM_pxbbq6V2zHyZCJmDk"
	DefaultAuthTokenKey = "#@!{[duXm-serVice-t0ken]},.(10086)$!"

	_, err := NewContext(nil)
	assert.Error(t, err)
}

func TestNewContextNoAuth(t *testing.T) {
	DefaultAuthToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJEdXhtLUp3dC1Ub2tlbiIsImV4cCI6MTcwMzc1Mjk3NCwidXNlciI6eyJ0ZW5hbnROYW1lIjoidGVzdCIsIm5hbWUiOiJ0ZXN0IiwidGVuYW50SWQiOiJ0ZXN0IiwidGVuYW50QWNjb3VudCI6InRlc3QiLCJpZCI6IjE3MjY0Nzk0NDEyNTUwNTEyNjQiLCJ1c2VyVHlwZSI6IlRFTkFOVF9BRE1JTiIsImFjY291bnQiOiJ0ZXN0Iiwic3RhdHVzIjoiVVNFSU5HIn0sImNsaWVudF9pZCI6IjA5OGY2YmNkNDYyMWQzNzNjYWRlNGU4MzI2MjdiNGY2In0.s_kHa3pKt6XehbsL7E9PJqywM_pxbbq6V2zHyZCJmDk"
	DefaultAuthTokenKey = "#@!{[duXm-serVice-t0ken]},.(10086)$!"

	_, err := NewContextNoAuth(nil)
	assert.NoError(t, err)
}
