package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func init() {
	Connection()
}
func Connection() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error to load the env file")
	}
	uri := os.Getenv("DB_URI")
	log.Println(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	log.Println("Connected to Database")

	return Client
	//return Client.Database("Dashboard").Collection("")
}
func OpenCollection(user string) *mongo.Collection {
	return Client.Database("Dashboard").Collection(user)
}
