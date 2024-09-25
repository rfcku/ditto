package vote

import (
	"context"
	"fmt"
	"go-api/internal/comment"
	"go-api/internal/post"
	"go-api/pkg/db"
	"time"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var voteCollection = db.GetCollection("votes")

func DbGetAllVotes(page string, limit string) ([]Vote, error) {
	var votes []Vote

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	fOpt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
	}

	cursor, err := voteCollection.Find(context.Background(), bson.M{}, &fOpt)

	if err != nil {
		return votes, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var vote Vote
		cursor.Decode(&vote)
		votes = append(votes, vote)
	}
	return votes, nil
}

func DbGetVoteID(id primitive.ObjectID) (Vote, error) {
	var vote Vote
	err := voteCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&vote)
	return vote, err
}

func DbVoteExists(postID primitive.ObjectID, user string) (bool, error) {
	count, err := voteCollection.CountDocuments(context.Background(), bson.M{"target_id": postID, "author_id": user})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return count > 0, nil
}


func DbCreateVote(vote Vote) (primitive.ObjectID, error) {
	result, err := voteCollection.InsertOne(context.Background(), vote)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateVote(id primitive.ObjectID, vote Vote) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": vote}
	_, err := voteCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteVote(id primitive.ObjectID) error {
	_, err := voteCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbDeleteVoteByAuthor(postID primitive.ObjectID, authorID string) error {
	_, err := voteCollection.DeleteOne(context.Background(), bson.M{"target_id": postID, "author_id": authorID})
	return err
}


func DbGetTargetVotes(postID primitive.ObjectID) (int32, error) {
	var votes []Vote
	cursor, err := voteCollection.Find(context.Background(), bson.M{"target_id": postID})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var vote Vote
		cursor.Decode(&vote)
		votes = append(votes, vote)
	}
	return int32(len(votes)), nil	
}


func DbSubmitVote(target string, targetID primitive.ObjectID, authorID string) (int32, bool, error) {

	if target == "p" {
		exists, err := post.DbPostExists(targetID)
		if err != nil {
			return 0, false, err
		}
		if !exists {
			return 0, false, fmt.Errorf("post not found")
		}	
	}

	if target == "c" {
		exists, err := comment.DbCommentExists(targetID)
		if err != nil {
			return 0, false, err
		}
		if !exists {
			return 0, false, fmt.Errorf("comment not found")
		}
	
	}

	var voted bool
	var v Vote
	v.TargetID = targetID
	v.AuthorID = authorID
	v.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	v.Type = target
	
	vote, err := DbVoteExists(v.TargetID, v.AuthorID)
	if err != nil {
		return 0, false, err
	}
	fmt.Println("Vote exists: ", vote)
	if vote {
		err = DbDeleteVoteByAuthor(v.TargetID, v.AuthorID)
		if err != nil {
			return 0, true, err
		}
		voted = false
	} else { 
		_,err := DbCreateVote(v)
		if err != nil {
			return 0, false, err
		}
		voted = true
	}
	votes, err := DbGetTargetVotes(v.TargetID)
	if err != nil {
		return 0, voted, err
	}
	return votes, voted, nil
}