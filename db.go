package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {

	// client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("mongo_db_url")))
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)

	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://unequaled-authority-0911:sbYJWxxyLqFujQVWwUkmYpDHq3jW8ZYfh5LxeQfE@andromeda-4pgy9.mongodb.net"))
	fmt.Println(os.Getenv("mongo_db_url"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("mongo_db_url")))
	err = client.Ping(ctx, readpref.Primary())
	// client, err := mongo.Connect(context.Background(), os.Getenv("mongo_db_url"), nil)
	return client, err
}

func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	// pipe := client.TxPipeline()
	// pipe.ZAdd("gliR:1", redis.Z{Score:235, Member: "shubh"})
	// pipe.ZAdd("gliR:1", redis.Z{Score:234, Member: "shubhnkar"})
	// pipe.Exec()
	return client, err
}
