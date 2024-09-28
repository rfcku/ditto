package post

import (
	"errors"
	utils "go-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p Post) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Content == "" {
		return errors.New("content is required")
	}
	if p.Link == "" {
		return errors.New("link is required")
	}
	if p.AuthorID == "" {
		return errors.New("author_id is required")
	}
	return nil
}

func (p Post) View() PostView {
	return PostView{
		ID:            utils.ObjectIdToString(p.ID),
		TargetID:      utils.ObjectIdToString(p.ID),
		Title:         p.Title,
		Content:       p.Content,
		AuthorID:      p.AuthorID,
		Link:          p.Link,
		Tags:          p.Tags,
		VotesTotal:    p.VotesTotal,
		Voted:         p.Voted,
		CommentsTotal: p.CommentsTotal,
		Awards:        p.Awards,
		AwardsTotal:   p.AwardsTotal,
		CreatedAt:     utils.DateToString(p.CreatedAt),
		Type:          p.Type,
	}
}

func ToPostView(posts []Post) []PostView {
	var postViews []PostView
	for _, post := range posts {
		postViews = append(postViews, post.View())
	}
	return postViews
}

func PostPipeline(pipeline []bson.M) []bson.M {
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

func PostsPaginatedPipeline(page int64, limit int64, sortBy string) []bson.M {
	skip := page*limit - limit

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
		}
		pipeline = PostPipeline(pipeline)
		return pipeline
	}
	pipeline := []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{"$skip": skip},
		{"$limit": limit},
	}
	pipeline = PostPipeline(pipeline)
	return pipeline
}

func PostByIDPipeline(postID primitive.ObjectID) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": postID}},
	}
	pipeline = PostPipeline(pipeline)
	return pipeline
}

func PostsByUserPipeline(username string) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"author_id": username}},
	}
	pipeline = PostPipeline(pipeline)
	return pipeline
}

func AppendSorter(pipeline []bson.M, sortBy string) []bson.M {
	if sortBy == "new" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "unvoted" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": -1}})
	}
	return pipeline
}

func AppendVotedToPipeline(pipeline []bson.M, authorID string) []bson.M {
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
