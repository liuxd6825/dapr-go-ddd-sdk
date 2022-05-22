package ddd_errors

import (
	"errors"
	"testing"
)

func TestIsGrpcConnError(t *testing.T) {
	err := errors.New("rpc error: code = Internal desc = fail to invoke, id: query-service, err: couldn't find service: query-service")
	if IsGrpcConnError(err) {
		t.Error("error ")
	}

	err = nil
	if IsGrpcConnError(err) {
		t.Error("error ")
	}

	err = errors.New("rpc error: code = Unavailable desc = connection error connect: connection refused")
	if !IsGrpcConnError(err) {
		t.Error("error ")
	}

}
