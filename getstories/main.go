package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/kelseyhightower/envconfig"
	wflambda "github.com/retgits/wavefront-lambda-go"
)

var wfAgent = wflambda.NewWavefrontAgent(&wflambda.WavefrontConfig{})

// config is the struct that is used to keep track of all environment variables
type config struct {
	ToAddress   string `required:"true" split_words:"true" envconfig:"TOADDRESS"`
	FromAddress string `required:"true" split_words:"true" envconfig:"FROMADDRESS"`
	AWSRegion   string `required:"true" split_words:"true" envconfig:"REGION"`
	APIToken    string `required:"true" split_words:"true" envconfig:"APITOKEN"`
	Days        int    `required:"true" split_words:"true" envconfig:"DAYS"`
	Owner       string `required:"true" split_words:"true" envconfig:"OWNER"`
}

var c config

// The handler function is executed every time that a new Lambda event is received.
// It takes a JSON payload (you can see an example in the event.json file) and only
// returns an error if the something went wrong. The event comes fom CloudWatch and
// is scheduled every interval (where the interval is defined as variable)
func handler(request events.CloudWatchEvent) error {
	// Get configuration set using environment variables
	err := envconfig.Process("", &c)
	if err != nil {
		log.Println(fmt.Sprintf("error starting function: %s", err.Error()))
		return err
	}

	searchInput := SearchStoriesByOwnerInput{
		APIToken: c.APIToken,
		Days:     c.Days,
		Owner:    c.Owner,
	}

	stories, err := SearchStoriesByOwner(searchInput)
	if err != nil {
		log.Println(fmt.Sprintf("error running function: %s", err.Error()))
		return err
	}

	// TODO: Replace this with a proper text template
	bodyContent := ""

	for _, story := range stories.Stories {
		deadline, err := time.Parse("2006-01-02T15:04:05Z", story.Deadline)
		if err != nil {
			fmt.Println("Error while parsing date :", err)
		}
		bodyContent = fmt.Sprintf("%s\n\n%s\nLink: %s\nDue on %s\n", bodyContent, story.Name, strings.ReplaceAll(story.AppURL, "\\", ""), deadline)
	}

	// Create an AWS session
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.AWSRegion),
	}))

	// Create an instance of the SES Session
	sesSession := ses.New(awsSession)

	// Create the Email request
	sesEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(c.ToAddress)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(bodyContent),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Your Clubhouse stories due in the next 7 days"),
			},
		},
		Source: aws.String(c.FromAddress),
		ReplyToAddresses: []*string{
			aws.String(c.FromAddress),
		},
	}

	// Send the email
	_, err = sesSession.SendEmail(sesEmailInput)
	return err
}

// The main method is executed by AWS Lambda and points to the handler
func main() {
	lambda.Start(wfAgent.WrapHandler(handler))
}
