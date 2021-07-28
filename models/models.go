package models

type FetchRequestModel struct {
	StartDate int     `json:"startDate" bson:"startDate"`
	EndDate   float64 `json:"endDate" bson:"endDate"`
	MinCount  string  `json:"minCount" bson:"minCount"`
	MaxCount  string  `json:"maxCount" bson:"maxCount"`
}
