package user

import (
	"context"
	"fmt"
	"go-api/pkg/db"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection = db.GetCollection("users")

func DbGetAllUsers(page string, limit string) ([]User, error) {
	var users []User

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	fOpt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
	}

	cursor, err := userCollection.Find(context.Background(), bson.M{}, &fOpt)

	if err != nil {
		return users, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	return users, nil
}

func DbGetUserByUsername(nickname string) (User, error) {
	var u User
	err := userCollection.FindOne(context.Background(), bson.M{"username": nickname}).Decode(&u)
	return u, err
}

func DbGetUserByID(id primitive.ObjectID) (User, error) {
	var user User
	err := userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func DbUserExists(username string) (bool, error) {
	count, err := userCollection.CountDocuments(context.Background(), bson.M{"username": username })
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return count > 0, nil
}


func DbCreateUser(user User) (primitive.ObjectID, error) {
	result, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateUser(id primitive.ObjectID, user User) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": user}
	_, err := userCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteUser(id primitive.ObjectID) error {
	_, err := userCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}