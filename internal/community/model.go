package community

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Community struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name 		string `json:"name" bson:"name"`
	Tags 		[]string `json:"tags" bson:"tags"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
}
