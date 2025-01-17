package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	start := time.Now()
	mongodbURI := "mongodb://localhost:27017"
	database := "attr"
	currentCollection := "attribute"
	batchSize := 10000

	newCollection := make(map[string]string)
	newCollection["device"] = "device_attribute"
	newCollection["user"] = "user_attribute"

	documents := make(map[string][]interface{})
	for key := range newCollection {
		documents[key] = []interface{}{}
	}

	clientOptions := options.Client().ApplyURI(mongodbURI)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	collection := client.Database(database).Collection(currentCollection)
	filter := bson.D{}
	findOptions := options.Find()
	findOptions.SetBatchSize(int32(batchSize))
	findOptions.SetNoCursorTimeout(true)

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for cursor.Next(context.Background()) {
		var record bson.M
		if err := cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}
		documents[record["entity_type"].(string)] = append(documents[record["entity_type"].(string)], record)
		count += 1

		if count%batchSize == 0 {
			for key, value := range newCollection {
				if len(documents[key]) == 0 {
					continue
				}
				_, err := client.Database(database).Collection(value).InsertMany(context.Background(), documents[key])
				if err != nil {
					log.Fatal("Failed to insert documents:", err)
				}
				documents[key] = []interface{}{}
			}
		}
	}

	for key, value := range newCollection {
		if len(documents[key]) == 0 {
			continue
		}
		_, err := client.Database(database).Collection(value).InsertMany(context.Background(), documents[key])
		if err != nil {
			log.Fatal("Failed to insert documents:", err)
		}
		documents[key] = []interface{}{}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total documents processed: %d\n", count)
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Total runtime: %v\n", duration)
}
