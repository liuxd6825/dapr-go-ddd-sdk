package ddd_errors

import "go.mongodb.org/mongo-driver/mongo"

func IsErrorMongoNoDocuments(err error) bool {
	if err == mongo.ErrNoDocuments {
		return true
	}
	return false
}
