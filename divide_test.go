package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestDivideCollection(t *testing.T) {
	//Connect to mongodb
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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
		log.Fatal("Could not connect to MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB")

	//Generate record and save to attribute collection
	generateRecord(client)

	//Get all record from attribute and divide into userDocuments slice and deviceDocuments slice
	deviceDocuments, userDocuments := getRecordFromAttribute(client)

	// Divide collection
	main()

	//Get all record from user_attribute and see if it equals generated user record slice
	collection := client.Database("attr").Collection("user_attribute")
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetBatchSize(2000)
	findOptions.SetNoCursorTimeout(true)
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	docs := []interface{}{}
	for cursor.Next(context.Background()) {
		var record bson.M
		if err := cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}
		docs = append(docs, record)
	}
	if !reflect.DeepEqual(userDocuments, docs) {
		t.Errorf("User documents not correctly divided")
	}

	//Get all record from device_attribute and see if it equals generated device record slice
	collection = client.Database("attr").Collection("device_attribute")

	cursor, err = collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	docs = []interface{}{}
	for cursor.Next(context.Background()) {
		var record bson.M
		if err := cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}
		docs = append(docs, record)
	}
	if !reflect.DeepEqual(deviceDocuments, docs) {
		t.Errorf("Device documents not correctly divided")
	}
}

func generateRecord(client *mongo.Client) {
	collection := client.Database("attr").Collection("attribute")

	var documents []interface{}
	for i := 1; i <= 1000; i++ {
		if i%2 == 0 {
			doc := bson.D{
				{"entity_id", fmt.Sprintf("Entity_%d", i)},
				{"entity_type", "device"},
				{"name", "abcd"},
				{"ab", 1},
				{"state", true},
				{"abcd", bson.D{
					{"a", "bc"},
					{"num", 123},
				}},
			}
			documents = append(documents, doc)
		} else {
			doc := bson.D{
				{"entity_id", fmt.Sprintf("Entity_%d", i)},
				{"entity_type", "device"},
			}
			documents = append(documents, doc)
		}
	}

	for i := 1; i <= 1000; i++ {
		if i%2 == 0 {
			doc := bson.D{
				{"entity_id", fmt.Sprintf("Entity_%d", i)},
				{"entity_type", "user"},
				{"name", "abcd"},
				{"ab", 1},
				{"state", true},
				{"abcd", bson.D{
					{"a", "bc"},
					{"num", 123},
				}},
			}
			documents = append(documents, doc)
		} else {
			doc := bson.D{
				{"entity_id", fmt.Sprintf("Entity_%d", i)},
				{"entity_type", "user"},
			}
			documents = append(documents, doc)
		}
	}

	result, err := collection.InsertMany(context.Background(), documents)
	if err != nil {
		log.Fatal("Failed to insert documents:", err)
	}

	fmt.Printf("Inserted %d documents\n", len(result.InsertedIDs))
}

func getRecordFromAttribute(client *mongo.Client) ([]interface{}, []interface{}) {
	collection := client.Database("attr").Collection("attribute")
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetBatchSize(2000)
	findOptions.SetNoCursorTimeout(true)

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	userDocuments := []interface{}{}
	deviceDocuments := []interface{}{}
	for cursor.Next(context.Background()) {
		var record bson.M
		if err := cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}
		if record["entity_type"].(string) == "device" {
			deviceDocuments = append(deviceDocuments, record)
		} else {
			userDocuments = append(userDocuments, record)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return deviceDocuments, userDocuments
}
