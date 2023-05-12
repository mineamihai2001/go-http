package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schema interface{}

type Category string

const (
	GENERAL   Category = "general"
	HISTORY   Category = "history"
	GEOGRAPHY Category = "geography"
	SPORT     Category = "sport"
	SCIENCE   Category = "science"
)

type Fact struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Category Category           `json:"category" bson:"category"`
	Question string             `json:"question" bson:"question"`
	Answer   string             `json:"answer" bson:"answer"`
}

type CreateFactDto struct {
	Category Category `json:"category"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
}

type UpdateFactDto struct {
	Id       string   `json:"id"`
	Category Category `json:"category"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
}

type DeleteFactDto struct {
	Category Category `json:"category"`
}

type CreateFactResponseDto struct {
	CreatedId string `json:"createdId"`
}

type DeleteFactResponseDto struct {
	DeletedCount int64 `json:"deletedCount"`
}

type UpdateFactResponseDto struct {
	UpdatedCount int64 `json:"updatedCount"`
}
