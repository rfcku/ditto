package post

import (
	"testing"
)

func Test_GetPipeline(t *testing.T) {
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	pipeline = GetPipeline(pipeline)
	if len(pipeline) != 2 {
		t.Errorf("Expected pipeline to have 2 stages, got %d", len(pipeline))
	}
}
