package comment

import (
	"go-api/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredFields(comment Comment) bool {
	if comment.TargetID == primitive.NilObjectID {
		return false
	}
	if comment.AuthorID == "" {
		return false

	}
	if comment.Content == "" {
		return false
	}
	return true
}

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func AddCommentsPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	if sortBy == "new" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "uncommentd" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"comments": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"comments": -1}})
	}
	return pipeline
}

func CommentsDefaultQueryParams(c *gin.Context) (int, int, string) {
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

func GetPipeline(pipeline []bson.M) []bson.M {
	dflt := []bson.M{
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "votes",
		}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "comments",
		}},
		{"$lookup": bson.M{
			"from":         "awards",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "awards",
		}},
		{"$addFields": bson.M{
			"votes_total": bson.M{"$size": "$votes"},
		}},
		{"$addFields": bson.M{
			"comments_total": bson.M{"$size": "$comments"},
		}},
		{"$addFields": bson.M{
			"awards_total": bson.M{"$size": "$awards"},
		}},
		{"$addFields": bson.M{
			"awards": bson.M{"$slice": []interface{}{"$awards", 3}},
		}},
	}
	pipeline = append(pipeline, dflt...)
	return pipeline
}

func GetCommentsPaginatedPipeline(page int, limit int, sortBy string, targetID primitive.ObjectID) []bson.M {
	skip := int64(page*limit - limit)
	if sortBy == "best" {
		pipeline := []bson.M{
			{
				"$match": bson.M{"target_id": targetID},
			},
			{
				"$lookup": bson.M{
					"from":         "votes",
					"localField":   "_id",
					"foreignField": "target_id",
					"as":           "votes",
				},
			},
			{
				"$addFields": bson.M{
					"votes_total": bson.M{"$size": "$votes"},
				},
			},
			{
				"$sort": bson.M{"votes_total": -1},
			},
			{
				"$skip": skip,
			},
			{
				"$limit": limit,
			},
		}
		pipeline = GetPipeline(pipeline)
		return pipeline
	}
	pipeline := []bson.M{
		{
			"$match": bson.M{"target_id": targetID},
		},
		{"$skip": skip},
		{"$limit": limit},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}

func GetCommentsPipeline(page int, limit int, sortBy string, user interface{}, targetID primitive.ObjectID) []bson.M {
	skip := int64(page*limit - limit)
	pipeline := []bson.M{
		{"$match": bson.M{"target_id": targetID}},
		{"$skip": skip},
		{"$limit": limit},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "votes",
		}},
		{"$addFields": bson.M{
			"votes_total": bson.M{"$size": "$votes"},
		}},
	}
	return pipeline
}

func GetPagination(page int, limit int, sortBy string, total int64, id primitive.ObjectID) Pagination {
	pagination := Pagination{}
	pagination.Page = page
	pagination.Limit = limit
	pagination.SortBy = sortBy
	pagination.TotalPages = total / int64(limit)
	pagination.TotalRecords = total
	pagination.CurrentPage = int64(page)
	if page < int(pagination.TotalPages) {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if page > 0 {
		pagination.HasPrev = true
	} else {
		pagination.HasPrev = false
	}
	pagination.NextLink = "/api/comments/" + ObjectIdToString(id) + "?page=" + strconv.Itoa(page+1) + "&limit=" + strconv.Itoa(limit)
	pagination.PrevLink = "/api/comments/" + ObjectIdToString(id) + "?page=" + strconv.Itoa(page-1) + "&limit=" + strconv.Itoa(limit)
	return pagination
}

func AddCommentsVotedPipeline(pipeline []bson.M, authorID string) []bson.M {
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from": "votes",
		"let":  bson.M{"target_id": "$_id"},
		"pipeline": []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$eq": []string{"$target_id", "$$target_id"}},
						{"$eq": []string{"$author_id", authorID}},
					},
				},
			}},
		},
		"as": "voted",
	}})
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"voted": bson.M{"$cond": bson.M{
				"if": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$voted"},
						0,
					},
				},
				"then": true,
				"else": false,
			}},
		},
	})
	return pipeline
}

func CommentToView(comment Comment) CommentView {
	commentView := CommentView{}
	commentView.ID = ObjectIdToString(comment.ID)
	commentView.Content = comment.Content
	commentView.TargetID = ObjectIdToString(comment.TargetID)
	commentView.AuthorID = comment.AuthorID
	commentView.CreatedAt = utils.DateToString(comment.CreatedAt)
	commentView.VotesTotal = comment.VotesTotal
	commentView.Voted = comment.Voted
	commentView.Replies = comment.Replies
	return commentView
}

