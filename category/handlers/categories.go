package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func ReptileHandler(db *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		val, err := db.LRange(ctx, "reptile", 0, -1).Result()
		if err != nil {
			logrus.Fatal(err)
			c.String(400, "an error occured", err.Error())
			return
		}
		c.JSON(200, gin.H{"data": val})
	}
}

func CatHandler(db *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		val, err := db.LRange(ctx, "cat", 0, -1).Result()
		if err != nil {
			logrus.Fatal(err)
			c.String(400, "an error occured", err.Error())
			return
		}
		c.JSON(200, gin.H{"data": val})
	}
}

func DogHandler(db *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		val, err := db.LRange(ctx, "dog", 0, -1).Result()
		if err != nil {
			logrus.Fatal(err)
			c.String(400, "an error occured", err.Error())
			return
		}
		c.JSON(200, gin.H{"data": val})
	}
}
