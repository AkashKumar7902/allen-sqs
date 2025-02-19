package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.mongodb.org/mongo-driver/bson"
)

var sqsConsumer *sqs.SQS

// ConsumeEvents listens to the SQS queue and processes events
func ConsumeEvents(sess *session.Session) {
	sqsConsumer = sqs.New(sess)
	for {
		// Receive messages from SQS
		msgs, err := sqsConsumer.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(20), // Long poll for 20 seconds
		})
		if err != nil {
			log.Println("Failed to receive messages:", err)
			continue
		}

		// Process each message
		for _, msg := range msgs.Messages {
			event := msg.Body
			log.Printf("Received event: %s", *event)

			// Insert event into MongoDB
			eventDocument := bson.D{{Key: "event", Value: *event}, {Key: "processed_at", Value: time.Now()}}
			_, err := col2.InsertOne(context.Background(), eventDocument)
			if err != nil {
				log.Println("Failed to insert event into MongoDB:", err)
			}

			// Delete the message from the queue
			_, err = sqsConsumer.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueUrl),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Println("Failed to delete message:", err)
			} else {
				log.Printf("Processed and deleted event: %s", *event)
			}
		}
	}
}
