package community

import (
	"errors"
	"go-api/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (c Community) Validate() error {
	if c.Name == "" {
		return errors.New("title is required")
	}
	return nil
}

func (c Community) View() CommunityView {
	return CommunityView{
		ID:        utils.ObjectIdToString(c.ID),
		Name:      c.Name,
		CreatedAt: utils.DateToString(c.CreatedAt),
	}
}

func ToCommunityView(communities []Community) []CommunityView {
	var communityViews []CommunityView
	for _, community := range communities {
		communityViews = append(communityViews, community.View())
	}
	return communityViews
}

func AddCommunitiesPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})

	return pipeline
}

func GetCommunitiesPipeline(page int64, limit int64, sortBy string) []bson.M {
	skip := page*limit - limit
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	return pipeline
}

