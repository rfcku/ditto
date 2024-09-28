package community

import (
	"context"
	"go-api/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var communityCollection = db.GetCollection("communities")
var communityTypeCollection = db.GetCollection("community_types")

func DbGetAllCommunities(page int64, limit int64, sortBy string, user interface{}) ([]Community, int64, error) {

	var communities []Community
	pipeline := GetCommunitiesPipeline(page, limit, sortBy)
	pipeline = AddCommunitiesPipelineSorter(pipeline, sortBy)
	cursor, err := communityCollection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return communities, 0, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var community Community
		cursor.Decode(&community)
		communities = append(communities, community)
	}

	total, err := communityCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return communities, 0, err
	}
	return communities, total, nil
}

func DbGetCommunityID(id primitive.ObjectID) (Community, error) {
	var community Community
	err := communityCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&community)
	return community, err
}
func DBGetCommunityByName(name string) (Community, error) {
	println("DBGetCommunityByName", name)
	var community Community
	err := communityCollection.FindOne(context.Background(), bson.M{"name": name}).Decode(&community)
	return community, err
}

func DbGetSearchCommunities(search string) ([]Community, error) {
	var communities []Community
	filter := bson.M{"name": bson.M{"$regex": search, "$options": "i"}}
	cursor, err := communityCollection.Find(context.Background(), filter)
	if err != nil {
		return communities, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var community Community
		cursor.Decode(&community)
		communities = append(communities, community)
	}
	return communities, nil
}

func DbCommunityTypeExists(id primitive.ObjectID) (bool, error) {
	count, err := communityTypeCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbCreateCommunity(community Community) (primitive.ObjectID, error) {
	result, err := communityCollection.InsertOne(context.Background(), community)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateCommunity(id primitive.ObjectID, community Community) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": community}
	_, err := communityCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteCommunity(id primitive.ObjectID) error {
	_, err := communityCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomCommunity() (Community, error) {
	var community Community
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}
	cursor, err := communityCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return community, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cursor.Decode(&community)
	}
	return community, nil
}
