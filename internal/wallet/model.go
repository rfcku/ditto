package wallet

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wallet struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Balance     int    `json:"balance" bson:"balance"`
	CreatedAt   string `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
	UpdatedAt   string `json:"updated_at" bson:"updated_at"`
}
