package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FetchRequestModel struct {
	StartDate string `json:"startDate" bson:"startDate"`
	EndDate   string `json:"endDate" bson:"endDate"`
	MinCount  int    `json:"minCount" bson:"minCount"`
	MaxCount  int    `json:"maxCount" bson:"maxCount"`
}

type FetchResponseModel struct {
	Code    int                      `json:"code" bson:"code"`
	Msg     string                   `json:"msg" bson:"msg"`
	Records []FetchRecordsArrayModel `json:"records" bson:"records"`
}

type FetchRecordsArrayModel struct {
	Key         string `json:"key" bson:"key"`
	CreatedAt   string `json:"createdAt" bson:"createdAt"`
	TotalCounts int    `json:"totalCounts" bson:"totalCounts"`
}

type DataFromMongoDBModel struct {
	Key         string             `json:"key" bson:"key"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	TotalCounts int                `json:"totalCounts" bson:"totalCounts"`
}

type PostRequestModel struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type PostResponseModel struct {
	Value string `json:"value" bson:"value"`
}

type GetResponseModel struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}
