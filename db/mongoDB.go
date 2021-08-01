package db

import (
	"context"
	"log"
	"os"
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

//Create a mongodb clients
func ConnectMongoDB() (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
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

//Fetches data from MongoDB
func (db *MongoDB) FetchDataFromMongoDB(startTime time.Time, endTime time.Time, minCount int, maxCount int) ([]models.FetchRecordsArrayModel, string, int) {

	database := db.MongoClient.Database("getir-case-study")
	recordsCollection := database.Collection("records")
	//The response time is short however I rather care about response than performance. I do not want my server to crash
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	//Here we basicly filter records considering their createdAt fields.
	matchStage := bson.D{{"$match", bson.D{primitive.E{Key: "createdAt", Value: bson.M{"$gt": startTime, "$lt": endTime}}}}}
	//we simply sum all the arrays. I placed importance on the case that there are multiple records with the same key.
	groupStage := bson.D{{"$project", bson.D{primitive.E{Key: "key", Value: "$key"}, primitive.E{Key: "createdAt", Value: "$createdAt"}, {"totalCounts", bson.D{{"$sum", "$counts"}}}}}}
	//we apply another filter here considering their array summations.
	matchStage2 := bson.D{{"$match", bson.D{primitive.E{Key: "totalCounts", Value: bson.M{"$gt": minCount, "$lt": maxCount}}}}}
	cursor, err := recordsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, matchStage2})
	if err != nil {
		log.Println(err)
		return nil, "pipeline did not work", 500
	}
	result := []models.DataFromMongoDBModel{}
	//we get all the data and set our array model
	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, "cursor has crushed", 500
	}
	response := make([]models.FetchRecordsArrayModel, len(result))
	for index, item := range result {
		response[index] = models.FetchRecordsArrayModel{item.Key, item.CreatedAt.Time().Format(time.RFC3339Nano), item.TotalCounts}
	}
	return response, "Success", 0
}
