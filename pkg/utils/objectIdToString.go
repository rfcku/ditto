package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}
