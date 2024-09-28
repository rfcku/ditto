package file

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path"
	"strconv"
)

func GetFileType(filename string) int8 {
	fileExtension := path.Ext(filename)
	var fileType int8
	if fileExtension == ".jpg" || fileExtension == ".jpeg" || fileExtension == ".png" {
		fileType = 1
	} else if fileExtension == ".mp4" || fileExtension == ".mov" {
		fileType = 2
	} else {
		fileType = 0
	}
	return fileType
}

func RequiredFields(file File) bool {
	if file.TargetID == primitive.NilObjectID {
		return false
	}

	if file.AuthorID == "" {
		return false
	}
	return true
}

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func DateToString(date primitive.DateTime) string {
	return date.Time().String()
}

func AddFilesPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	if sortBy == "new" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "unfiled" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"files": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"files": -1}})
	}
	return pipeline
}

func FilesDefaultQueryParams(c *gin.Context) (int, int, string) {
	page := c.Query("page")
	limit := c.Query("limit")
	sortBy := c.Query("sort_by")
	var p, l int = 0, 0
	if page == "" {
		p = 1
	} else {
		p, _ = strconv.Atoi(page)
	}
	if limit == "" {
		l = 10
	} else {
		l, _ = strconv.Atoi(limit)
	}
	if sortBy == "" {
		sortBy = "best"
	}

	if l > 100 {
		l = 100
	}
	return p, l, sortBy
}

func GetFilesPipeline(page int, limit int, sortBy string) []bson.M {
	skip := int64(page*limit - limit)
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	return pipeline
}

