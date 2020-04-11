package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
)

// ConnectMongo creates a mongo connection
// func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("mongo_db_url")))
// 	err = client.Ping(ctx, readpref.Primary())

// 	return client, err
// }

// ConnectRedis creates a redis connection
func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client, err
}

// ConnectSQL creates a MySQL connection
func ConnectSQL() (*gorm.DB, error) {
	host := os.Getenv("sequelize_host")
	port := os.Getenv("sequelize_port")
	username := os.Getenv("sequelize_user")
	password := os.Getenv("sequelize_pass")
	database := os.Getenv("sequelize_db")
	dbSource := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True",
		username,
		password,
		host,
		port,
		database,
	)
	db, err := gorm.Open("mysql", dbSource)
	//db.AutoMigrate(&User{})
	// db.LogMode(true)
	return db, err
}
