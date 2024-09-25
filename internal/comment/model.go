package comment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	Content 	string `json:"content" bson:"content"`
	Replies		CommentPaginated `json:"replies" bson:"replies,omitempty"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at" swaggertype:"primitive,string"`
	Voted 		bool `json:"voted" bson:"voted,omitempty"`
	VotesTotal int32 `json:"votes_total" bson:"votes_total,omitempty"`
}

type CommentView struct {
	ID      	string `json:"id"`
	TargetID 	string `json:"target_id"`
	AuthorID 	string `json:"author_id"`
	Content 	string `json:"content"`
	Replies		CommentPaginated `json:"replies"`
	CreatedAt  string `json:"created_at"`
	Voted		bool `json:"voted"`
	VotesTotal int32 `json:"votes_total"`
}

type Pagination struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	SortBy string `json:"sort_by"`
	HasPrev bool `json:"has_prev"`
	HasNext bool `json:"has_next"`
	TotalRecords int64 `json:"total_records"`
	TotalPages int64 `json:"total_pages"`
	CurrentPage int64 `json:"current_page"`
	NextLink string `json:"next_link"`
	PrevLink string `json:"prev_link"`
}

type CommentPaginated struct {
	Data []CommentView `json:"data"`
	Pagination Pagination `json:"pag"`
}
