package feign_pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewURLParser(t *testing.T) {
	p, err := NewURLParser("dapr://user-service/api/v1/tenants/test/users")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, p.ServiceName, "user-service")
	assert.Equal(t, p.Path, "/api/v1/tenants/test/users")
	assert.Equal(t, p.Protocol.String(), "dapr")

	p, err = NewURLParser("http://user-service:8080/api/v1/tenants/test/users")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, p.ServiceName, "user-service")
	assert.Equal(t, p.Port, "8080")
	assert.Equal(t, p.Path, "/api/v1/tenants/test/users")
	assert.Equal(t, p.Protocol.String(), "http")
}
