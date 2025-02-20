// Package main handle the server
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var col *mongo.Collection
var col2 *mongo.Collection
var logger *zap.Logger

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using environment variables.")
	}

	time.Sleep(2 * time.Second)
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if err != nil {
			log.Println(err)
		}
	}()

	dbName, collection := "keploy", "url-shortener"
	collection2 := "events"

	client, err := New("localhost:27017", dbName)
	if err != nil {
		logger.Fatal("failed to create mgo db client", zap.Error(err))
	}
	db := client.Database(dbName)

	col = db.Collection(collection)
	col2 = db.Collection(collection2)

	port := "8080"

	println("PID:", os.Getpid())

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		log.Fatal("Missing AWS credentials in environment variables")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Fatal("Failed to create AWS session:", err)
	}

	// Start AppConfig fetcher
	go FetchAppConfig(sess)

	// Start SQS producer (optional)
	go SendEvent(sess)

	// Start SQS consumer
	go ConsumeEvents(sess)

	r := gin.Default()

	r.GET("/:param", getURL)
	r.POST("/url", putURL)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopper

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server Shutdown:", zap.Error(err))
		}
		logger.Info("stopper called")
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("listen: %s\n", zap.Error(err))
	}
	logger.Info("server exiting")
}
