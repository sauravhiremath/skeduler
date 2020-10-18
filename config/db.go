package config

import (
	"context"
	"log"
	"os"
	"skeduler/controllers"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbname string = os.Getenv("mongo_dbname")
var password string = os.Getenv("mongo_password")

// Connect initialises mongoDB driver
func Connect() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://saurav:" + password + "@skedules.tzfcm.gcp.mongodb.net/" + dbname + "?retryWrites=true&w=majority")
	client, err := mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	//Cancel context to avoid memory leak
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
	db := client.Database("go_mongo")
	controllers.MeetingCollection(db)
	// controllers.ParticipantCollection(db);
	return
}
