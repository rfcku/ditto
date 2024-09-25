package community

import (
	"context"
	"go-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var awardCollection = db.GetCollection("awards")
var awardTypeCollection = db.GetCollection("award_types")

func DbGetAllCommunities(page int, limit int, sortBy string, user interface{} ) ([]Community, error) {

	var awards []Community
	pipeline := GetCommunitiesPipeline(page, limit, sortBy)
	pipeline = AddCommunitiesPipelineSorter(pipeline, sortBy)

	cursor, err := awardCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return awards, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var award Community
		cursor.Decode(&award)
		awards = append(awards, award)
	}

	return awards, nil
}

func DbGetCommunityID(id primitive.ObjectID) (Community, error) {
	var award Community
	err := awardCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&award)
	return award, err
}

func DbCommunityTypeExists(id primitive.ObjectID) (bool, error) {
	count, err := awardTypeCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbCreateCommunity(award Community) (primitive.ObjectID, error) {
	result, err := awardCollection.InsertOne(context.Background(), award)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateCommunity(id primitive.ObjectID, award Community) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": award}
	_, err := awardCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteCommunity(id primitive.ObjectID) error {
	_, err := awardCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomCommunity() (Community, error) {
	var award Community
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
