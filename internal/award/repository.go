package award

import (
	"context"
	"go-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var awardCollection = db.GetCollection("awards")
var awardTypeCollection = db.GetCollection("award_types")

func DbGetAllAwards(page int, limit int, sortBy string, user interface{} ) ([]Award, error) {

	var awards []Award
	pipeline := GetAwardsPipeline(page, limit, sortBy)
	pipeline = AddAwardsPipelineSorter(pipeline, sortBy)

	cursor, err := awardCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return awards, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var award Award
		cursor.Decode(&award)
		awards = append(awards, award)
	}

	return awards, nil
}

func DbGetAwardID(id primitive.ObjectID) (Award, error) {
	var award Award
	err := awardCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&award)
	return award, err
}

func DbAwardTypeExists(id primitive.ObjectID) (bool, error) {
	count, err := awardTypeCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbCreateAward(award Award) (primitive.ObjectID, error) {
	result, err := awardCollection.InsertOne(context.Background(), award)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateAward(id primitive.ObjectID, award Award) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": award}
	_, err := awardCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteAward(id primitive.ObjectID) error {
	_, err := awardCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomAward() (Award, error) {
	var award Award
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}
	cursor, err := awardCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return award, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cursor.Decode(&award)
	}
	return award, nil
}
