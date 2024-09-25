package vote

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Vote struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
	Type 		string `json:"type" bson:"type"`
}
