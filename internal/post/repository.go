package post

import (
	"context"
	"go-api/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var postCollection = db.GetCollection("posts")

func DbGetAllPosts(page int, limit int, sortBy string, user interface{} ) ([]Post, int64, error) {

	var posts []Post
	pipeline := GetPostsPaginatedPipeline(page, limit, sortBy)
	
	if user != nil {
		nickname := user.(map[string]interface{})["nickname"].(string)
		if nickname != "" {
			pipeline = AddPostsVotedPipeline(pipeline, nickname)
		}
	}
	pipeline = AddPostsPipelineSorter(pipeline, sortBy)

	cursor, err := postCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return posts, 0, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	totalRecords, err := postCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return posts, 0, err
	}

	return posts, totalRecords, nil
}

func DbGetPostsByUser(username string) ([]Post, error) {
	var posts []Post
	pipeline := GetPostsByUserPipeline(username)
	cursor, err := postCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return posts, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}

	return posts, nil
}

func DbGetPostID(id primitive.ObjectID, user interface{}) (Post, error) {
	
	pipeline := GetPostPipelineByID(id)
	if user != nil {
		nickname := user.(map[string]interface{})["nickname"].(string)
		pipeline = AddPostsVotedPipeline(pipeline, nickname)
	}

	var post Post
	cursor, err := postCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return post, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cursor.Decode(&post)
	}
	
	return post, nil
}

func DbPostExists(id primitive.ObjectID) (bool, error) {
	count, err := postCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbCreatePost(post Post) (primitive.ObjectID, error) {
	post.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	result, err := postCollection.InsertOne(context.Background(), post)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdatePost(id primitive.ObjectID, post Post) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": post}
	_, err := postCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeletePost(id primitive.ObjectID) error {
	_, err := postCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomPost() (Post, error) {
	var post Post
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	cursor, err := postCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return post, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cursor.Decode(&post)
	}

	return post, nil
}