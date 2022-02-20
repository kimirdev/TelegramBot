package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoUrlPattern = "mongodb+srv://admin:%s@cluster0.khbsd.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

type MongoDB struct {
	database *mongo.Database
	client   *mongo.Client
}

type Link struct {
	Title string `bson:"title"`
	Url   string `bson:"url"`
}

func createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func NewMongoDB() (*MongoDB, error) {

	mongoPassword, exists := os.LookupEnv("MONGO_PASSWORD")

	if !exists {
		return nil, errors.New("MONGO_PASSWORD not found in file .env")
	}

	mongoUrl := fmt.Sprintf(mongoUrlPattern, mongoPassword)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}

	ctx, cf := createCtx()
	defer cf()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := client.Database("telegrambot")

	db.Collection("links")

	return &MongoDB{database: db, client: client}, nil
}

func (db *MongoDB) InsertLink(chatId int64, title, url string) error {
	linksCollection := db.database.Collection("links")
	ctx, cf := createCtx()
	defer cf()

	_, err := linksCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: title},
		{Key: "url", Value: url},
		{Key: "chatId", Value: chatId},
	})

	return err
}

func (db *MongoDB) GetAllLinks(chatId int64) ([]Link, error) {
	links := make([]Link, 0)

	ctx, cf := createCtx()
	defer cf()

	linksCollection := db.database.Collection("links")

	cursor, err := linksCollection.Find(ctx, bson.M{"chatId": chatId})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var link Link
		if err = cursor.Decode(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (db *MongoDB) DeleteLink(chatId int64, title, url string) error {
	linksCollection := db.database.Collection("links")

	ctx, cf := createCtx()
	defer cf()

	_, err := linksCollection.DeleteOne(ctx, bson.M{
		"title":  title,
		"url":    url,
		"chatId": chatId,
	})

	return err
}

func (db *MongoDB) Disconnect() {
	ctx, cf := createCtx()
	defer cf()
	db.client.Disconnect(ctx)
}
