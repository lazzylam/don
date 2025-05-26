package database

import (
    "anti-gcast-bot/config"
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
    client     *mongo.Client
    db         *mongo.Database
    ctx        = context.TODO()
)

func InitDB(cfg *config.Config) {
    var err error
    
    clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Ping database
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

    db = client.Database(cfg.DBName)
    log.Println("Connected to MongoDB!")
}

func GetCollection(name string) *mongo.Collection {
    return db.Collection(name)
}

func CloseDB() {
    if err := client.Disconnect(ctx); err != nil {
        log.Fatal(err)
    }
}