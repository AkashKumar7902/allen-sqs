package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsClient *sqs.SQS
var queueUrl = "https://sqs.eu-north-1.amazonaws.com/932089138525/allen-test"

// SendEvent produces random events and sends them to the SQS queue
func SendEvent(sess *session.Session) {
	sqsClient = sqs.New(sess)

	for {
		event := fmt.Sprintf("RandomEvent-%d", rand.Int())
		_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(event),
			QueueUrl:    aws.String(queueUrl),
		})
		if err != nil {
			log.Println("Failed to send message:", err)
		} else {
			log.Printf("Sent event: %s", event)
		}
		time.Sleep(time.Duration(40) * time.Second)
	}
}
