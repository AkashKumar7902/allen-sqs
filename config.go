package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
)

var appConfigClient *appconfig.AppConfig

// FetchAppConfig periodically fetches configuration from AWS AppConfig every 10-20 seconds.
func FetchAppConfig(sess *session.Session) {
	appConfigClient = appconfig.New(sess)

	// Define your AppConfig details
	applicationID := "meolhvk"
	environmentID := "testing"
	configurationProfileID := "m0sw0id"

	input := &appconfig.GetConfigurationInput{
		Application:   aws.String(applicationID),
		Environment:   aws.String(environmentID),
		Configuration: aws.String(configurationProfileID),
		ClientId:      aws.String("keploy"), // Unique ID for the client
	}

	output, err := appConfigClient.GetConfiguration(input)
	if err != nil {
		log.Printf("Failed to fetch configuration: %v", err)
		return
	}

	fmt.Printf("Fetched configuration: %s\n", string(output.Content))

	// Create a ticker that triggers every 10 seconds
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			output, err := appConfigClient.GetConfiguration(input)
			if err != nil {
				log.Printf("Failed to fetch configuration: %v", err)
				continue
			}

			fmt.Printf("Fetched configuration: %s\n", string(output.Content))

			// Reset the ticker to a new random duration between 10-20 seconds
			ticker.Reset(time.Duration(10) * time.Second)
		}
	}
}
