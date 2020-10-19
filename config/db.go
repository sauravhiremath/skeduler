package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/sauravhiremath/skeduler/controllers"
)

var dbname string = os.Getenv("MONGO_DBNAME")
var password string = os.Getenv("MONGO_PASSWORD")

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
		log.Fatal("[x] Couldn't connect to the database", err)
	} else {
		log.Println("[*] Connected!")
	}
	skedulerDB := client.Database("skeduler")
	controllers.Collection(skedulerDB)

	return
}
