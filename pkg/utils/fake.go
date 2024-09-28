package utils

import (
	"github.com/go-faker/faker/v4"
)

func FakePost() string {
	return `{
		"title": "` + faker.Word() + `",
		"content": "` + faker.Sentence() + `",
		"link": "` + faker.URL() + `",
		"author_id": "` + faker.Username() + `",
		"tags": ["` + faker.Word() + `", "` + faker.Word() + `"]
	}`
}
