package post

import (
	aw "go-api/internal/award"
	cm "go-api/internal/comment"
	utils "go-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Validator interface {
	Validate() error
}

type Post struct {
	ID            primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Title         string                `json:"title" bson:"title"`
	Content       string                `json:"content" bson:"content"`
	AuthorID      string                `json:"author_id" bson:"author_id,omitempty"`
	Link          string                `json:"link" bson:"link"`
	Tags          []string              `json:"tags" bson:"tags,omitempty"`
	Community     primitive.ObjectID    `json:"community" bson:"community"`
	VotesTotal    int32                 `json:"votes_total" bson:"votes_total,omitempty"`
	Voted         bool                  `json:"voted" bson:"voted,omitempty"`
	CreatedAt     primitive.DateTime    `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
	CommentsTotal int                   `json:"comments_total" bson:"comments_total,omitempty"`
	Comments      []cm.CommentPaginated `json:"comments" bson:"comments,omitempty"`
	Image         string                `json:"image" bson:"image,omitempty"`
	Awards        []aw.Award            `json:"awards" bson:"awards,omitempty"`
	AwardsTotal   int                   `json:"awards_total" bson:"awards_total,omitempty"`
	Files         []primitive.ObjectID  `json:"files" bson:"files,omitempty"`
	Type          int8                  `json:"type" bson:"type"`
}

type PostView struct {
	ID            string     `json:"id" bson:"_id,omitempty"`
	Title         string     `json:"title" bson:"title"`
	Content       string     `json:"content" bson:"content"`
	AuthorID      string     `json:"author_id" bson:"author_id,omitempty"`
	Link          string     `json:"link" bson:"link"`
	Tags          []string   `json:"tags" bson:"tags,omitempty"`
	Community     string     `json:"community" bson:"community"`
	Image         string     `json:"image" bson:"image,omitempty"`
	VotesTotal    int32      `json:"votes_total" bson:"votes_total,omitempty"`
	Voted         bool       `json:"voted" bson:"voted,omitempty"`
	CommentsTotal int        `json:"comments_total" bson:"comments_total,omitempty"`
	Awards        []aw.Award `json:"awards" bson:"awards,omitempty"`
	AwardsTotal   int        `json:"awards_total" bson:"awards_total,omitempty"`
	CreatedAt     string     `json:"created_at" bson:"created_at"`
	TargetID      string     `json:"target_id" bson:"target_id"`
	Files         []string   `json:"files" bson:"files,omitempty"`
	Type          int8       `json:"type" bson:"type"`
}

type PostPaginated struct {
	Data       []PostView       `json:"data"`
	Pagination utils.Pagination `json:"pagination"`
}

type PostIncoming struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Link    string `json:"link" bson:"link"`
	Tags    string `json:"tags" bson:"tags"`
}
