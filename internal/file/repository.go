package file

import (
	"context"
	"go-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var fileCollection = db.GetCollection("files")

func DbGetAllFiles(page int, limit int, sortBy string, user interface{} ) ([]File, error) {

	var files []File
	pipeline := GetFilesPipeline(page, limit, sortBy)
	pipeline = AddFilesPipelineSorter(pipeline, sortBy)

	cursor, err := fileCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return files, err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var file File
		cursor.Decode(&file)
		files = append(files, file)
	}

	return files, nil
}

func DbGetFileID(id primitive.ObjectID) (File, error) {
	var file File
	err := fileCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&file)
	return file, err
}

func DbCreateFile(file File) (primitive.ObjectID, error) {
	result, err := fileCollection.InsertOne(context.Background(), file)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateFile(id primitive.ObjectID, file File) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": file}
	_, err := fileCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteFile(id primitive.ObjectID) error {
	_, err := fileCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
