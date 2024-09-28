package file

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type File struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID  primitive.ObjectID `json:"target_id" bson:"target_id"`
	Type      int8               `json:"type" bson:"type"`
	AuthorID  string             `json:"author_id" bson:"author_id"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
}

type FileView struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	TargetID  string `json:"target_id" bson:"target_id"`
	TypeID    string `json:"type" bson:"type"`
	AuthorID  string `json:"author_id" bson:"author_id"`
	CreatedAt string `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
}

type Files []File
