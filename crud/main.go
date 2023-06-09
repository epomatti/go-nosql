package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(">> Disconnected!")
		}
	}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(">> Connected!")

	// Insert One

	type Actor struct {
		FirstName string
		LastName  string
		Awards    int16
	}

	collection := client.Database("dvdstore").Collection("actordetails")

	james := Actor{"James", "Roger", 9}

	insertResult, err := collection.InsertOne(context.TODO(), james)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a new actor: ", insertResult.InsertedID)

	// Insert Many
	mark := Actor{"Mark", "Brown", 0}
	mili := Actor{"Mili", "Ford", 1}

	actors := []interface{}{mark, mili}

	insertManyResult, err := collection.InsertMany(context.TODO(), actors)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple actors: ", insertManyResult.InsertedIDs)

	// Find
	// filter := bson.D{}
	filter := bson.D{{Key: "firstname", Value: "Mili"}}
	var result Actor

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found actor: %+v\n", result)

	// Update
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "awards", Value: 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v actors and updated %v actors.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Delete
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v actors\n", deleteResult.DeletedCount)
}
