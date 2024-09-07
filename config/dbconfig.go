package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

const dbName = "goBlogApp"

func ConnectDb() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	appEnv := os.Getenv("APP_ENV")
	uri := ""

	if appEnv == "production" {
		uri = os.Getenv("MONGO_URI_PUBLIC")
	} else {
		uri = os.Getenv("MONGO_URI_LOCAL")
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	Database = client.Database(dbName)
	fmt.Println("DB connected")
	return nil
}
