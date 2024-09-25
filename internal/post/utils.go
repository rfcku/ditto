package post

import (
	utils "go-api/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredFields(post Post) bool {
	if post.Title == "" {
		return false
	}
	if post.Content == "" {
		return false
	}
	if post.Link == "" {
		return false
	}
	if post.AuthorID == "" {
		return false
	}
	return true
}

func fakePost() string {
	return `{
		"title": "`+faker.Word()+`",
		"content": "`+faker.Sentence()+`",
		"link": "`+faker.URL()+`",
		"author_id": "`+faker.Username()+`",
		"tags": ["`+faker.Word()+`", "`+faker.Word()+`"]
	}`
}

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func PostToPostView(post Post) PostView {
	return PostView{
		ID: ObjectIdToString(post.ID),
		TargetID: ObjectIdToString(post.ID),
		Title: post.Title,
		Content: post.Content,
		AuthorID: post.AuthorID,
		Link: post.Link,
		Tags: post.Tags,
		VotesTotal: post.VotesTotal,
		Voted: post.Voted,
		CommentsTotal: post.CommentsTotal,
		Awards: post.Awards,
		AwardsTotal: post.AwardsTotal,
		CreatedAt: utils.DateToString(post.CreatedAt),
	}
}

func PostsToPostView(posts []Post) []PostView {
	var postView []PostView
	for _, post := range posts {
		postView = append(postView, PostToPostView(post))
	}
	return postView
}

func GetPagination(page int, limit int, sortBy string, total int64) Pagination {
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

	if page == 0 {
		pagination.HasPrev = false
	} else {
		pagination.HasPrev = true
	}
	pagination.NextLink = "?page=" + strconv.Itoa(page+1) + "&limit=" + strconv.Itoa(limit)
	pagination.PrevLink = "?page=" + strconv.Itoa(page-1) + "&limit=" + strconv.Itoa(limit)
	return pagination
}

func PostPaginatedView(posts []Post, totalRecords int64, page int, limit int, sortBy string) PostPaginated {
	PostPaginated := PostPaginated{}
	PostPaginated.Data = PostsToPostView(posts)

	pagination := Pagination{}
	pagination.Page = page
	pagination.Limit = limit
	pagination.SortBy = sortBy
	pagination.TotalPages = totalRecords / int64(limit)
	pagination.TotalRecords = totalRecords
	pagination.CurrentPage = int64(page)
	
	if page < int(pagination.TotalPages) {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if page == 0 {
		pagination.HasPrev = false
	} else {
		pagination.HasPrev = true
	}
	pagination.NextLink = "?page=" + strconv.Itoa(page+1) + "&limit=" + strconv.Itoa(limit)
	pagination.PrevLink = "?page=" + strconv.Itoa(page-1) + "&limit=" + strconv.Itoa(limit)
	PostPaginated.Pagination = pagination

	return PostPaginated
}

func AddPostsPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	if sortBy == "new"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "unvoted"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": -1}})
	}
	return pipeline
}

func PostsDefaultQueryParams(c *gin.Context) (int, int, string) {
	page := c.Query("page")
	limit := c.Query("limit")
	sortBy := c.Query("sort_by")
	var p, l int = 0, 0
	if page == "" {
		p = 1
	} else {
		p,_ = strconv.Atoi(page)
	}
	if limit == "" {
		l = 10
	} else {
		l,_ = strconv.Atoi(limit)
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
			// total awards and first 5 awards
			"awards": bson.M{"$slice": []interface{}{"$awards", 3}},
		}},
	}
	pipeline = append(pipeline, dflt...)
	return pipeline
}

func GetPostPipelineByID(postID primitive.ObjectID) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": postID}},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}
func GetPostsByUserPipeline(username string) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"author_id": username}},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}

func GetPostsPaginatedPipeline(page int, limit int, sortBy string) []bson.M {
	skip := int64(page*limit - limit)
	
	// lookup most voted posts
	if sortBy == "best" {
		pipeline := []bson.M{
			{
				"$addFields": bson.M{
					"totalRecords": "$totalRecords.count",
				},
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
		};
		pipeline = GetPipeline(pipeline)
		return pipeline
	}
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}

func AddPostsVotedPipeline(pipeline []bson.M, authorID string) []bson.M {
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "votes",
		"let":          bson.M{"target_id": "$_id"},
		"pipeline":     []bson.M{
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
				"if":  bson.M{
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
