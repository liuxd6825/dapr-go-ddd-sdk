package errors

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const (
	GrpcConnErrorPrefix = "rpc error: code = Unavailable desc = connection error"
	GrpcConnErrorSuffix = "connect: connection refused"
)

func IsErrorMongoNoDocuments(err error) bool {
	if err == mongo.ErrNoDocuments {
		return true
	}
	return false
}

//
// IsGrpcConnError
// @Description: 是否是GRPC连接错误
// @param err
// @return bool
//
func IsGrpcConnError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.HasPrefix(msg, GrpcConnErrorPrefix) {
		return true
	}

	return false
}
