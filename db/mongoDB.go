package db

import (
	"context"
	"fmt"
	"log"
	"time"

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

func (db *MongoDB) FetchDataFromMongoDB() {

	database := db.MongoClient.Database("getir-case-study")
	recordsCollection := database.Collection("records")

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	fmt.Println("111")

	//var dataFromMongoDB []models.DataFromMongoDBModel

	startDate := "2011-01-26"
	endDate := "2022-02-02"
	const (
		layoutISO = "2006-01-02"
	)

	t2, _ := time.Parse(layoutISO, startDate)

	t3, _ := time.Parse(layoutISO, endDate)

	fmt.Println(primitive.NewDateTimeFromTime(t2))

	fmt.Println(primitive.NewDateTimeFromTime(t3))

	matchStage := bson.D{{"$match", bson.D{primitive.E{Key: "key", Value: "TAKwGc6Jr4i8Z487"}, primitive.E{Key: "createdAt", Value: bson.M{"$gt": t2, "$lt": t3}}}}}
	groupStage := bson.D{{"$project", bson.D{{"key", "$key"}, {"createdAt", "$createdAt"}, {"totalCounts", bson.D{{"$sum", "$counts"}}}}}}

	cursor, err := recordsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		fmt.Println("12121")
		panic(err)
	}
	result := []bson.M{}

	if err = cursor.All(ctx, &result); err != nil {
		fmt.Println("333")
		panic(err)
	}
	fmt.Println(result)

}

func HelloMongo() {
	fmt.Println("Hello from Mongo")
}
