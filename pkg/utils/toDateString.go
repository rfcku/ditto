package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toDateString(date primitive.DateTime) string {
	return date.Time().String()
}
