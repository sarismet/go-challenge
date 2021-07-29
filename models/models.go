package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FetchRequestModel struct {
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
	MinCount  int    `json:"minCount" bson:"minCount"`
	MaxCount  int    `json:"maxCount" bson:"maxCount"`
}

type DataFromMongoDBModel struct {
	Key       string             `json:"key" bson:"key"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Counts    []int              `json:"counts" bson:"counts"`
	Value     string             `json:"value" bson:"value"`
}
