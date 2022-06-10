package ddd_mongodb

import "github.com/google/uuid"

func GetObjectID(id string) (ObjectId, error) {
	return ObjectId(id), nil
}

func NewObjectID() ObjectId {
	id := uuid.New().String()
	return ObjectId(id)
}
