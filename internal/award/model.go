package award

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AwardType struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name 		string `json:"name" bson:"name"`
}

type Award struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	TypeID 		primitive.ObjectID `json:"type" bson:"type"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
}

type AwardView struct {
	ID      	string `json:"id" bson:"_id,omitempty"`
	TargetID 	string `json:"target_id" bson:"target_id"`
	TypeID 		string `json:"type" bson:"type"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
}
