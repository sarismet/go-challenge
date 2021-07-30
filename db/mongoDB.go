package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-challenge/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDB struct {
	MongoClient *mongo.Client
}

func ConnectMongoDB() (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &MongoDB{
		MongoClient: client,
	}, nil
}
func (db *MongoDB) FetchDataFromMongoDB(startDate string, endDate string, minCount int, maxCount int) ([]models.FetchRecordsArrayModel, string, int) {

	database := db.MongoClient.Database("getir-case-study")
	recordsCollection := database.Collection("records")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	const (
		layoutISO = "2006-01-02"
	)

	startTime, _ := time.Parse(layoutISO, startDate)

	endTime, _ := time.Parse(layoutISO, endDate)

	matchStage := bson.D{{"$match", bson.D{primitive.E{Key: "createdAt", Value: bson.M{"$gt": startTime, "$lt": endTime}}}}}
	groupStage := bson.D{{"$project", bson.D{primitive.E{Key: "key", Value: "$key"}, primitive.E{Key: "createdAt", Value: "$createdAt"}, {"totalCounts", bson.D{{"$sum", "$counts"}}}}}}
	matchStage2 := bson.D{{"$match", bson.D{primitive.E{Key: "totalCounts", Value: bson.M{"$gt": minCount, "$lt": maxCount}}}}}

	cursor, err := recordsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, matchStage2})
	if err != nil {
		log.Println(err)
		return nil, "pipeline did not work", 500
	}

	result := []models.DataFromMongoDBModel{}
	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, "cursor has crushed", 500
	}

	fmt.Println(result)

	response := make([]models.FetchRecordsArrayModel, len(result))
	for index, item := range result {
		response[index] = models.FetchRecordsArrayModel{item.Key, item.CreatedAt.Time().Format(time.RFC3339Nano), item.TotalCounts}
	}

	return response, "Success", 0

}

func HelloMongo() {
	fmt.Println("Hello from Mongo")
}
