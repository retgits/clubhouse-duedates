# creditcard

[![Go Report Card](https://goreportcard.com/badge/github.com/retgits/clubhouse-duedates)](https://goreportcard.com/report/github.com/retgits/clubhouse-duedates)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/retgits/clubhouse-duedates)
![GitHub](https://img.shields.io/github/license/retgits/clubhouse-duedates)

> A serverless app to alert you on all your upcoming deadlines from Clubhouse.

## Pre-requisites

* [Go (at least Go 1.12)](https://golang.org/dl/)
* [A Wavefront API token](https://wavefront.com)
* [A Clubhouse API token](https://help.clubhouse.io/hc/en-us/articles/205701199-Clubhouse-API-Tokens)
* [An AWS account with access to SES](https://aws.amazon.com/ses/)

## Installation

### Get sources

To install you can clone this repository, and run `make deps` to get the [Go modules](./go.mod) this app relies on.

### Using make

There are a bunch of Makefile targets that help build and deploy the app

| Target  | Description                                                |
|---------|------------------------------------------------------------|
| build   | Build the executable for Lambda                            |
| clean   | Remove all generated files                                 |
| deploy  | Deploy the app to AWS Lambda                               |
| deps    | Get the Go modules from the GOPROXY                        |
| destroy | Deletes the CloudFormation stack and all created resources |
| help    | Displays the help for each target (this message)           |
| local   | Run SAM to test the Lambda function using Docker           |
| test    | Run all unit tests and print coverage                      |

## Configuration

Inside the [template.yaml](./template.yaml), there are a few configuration options that you can set:

```yaml
Environment:
    Variables:
        REGION: us-west-2 ## The AWS region you want to deploy your app to
        TOADDRESS: !Ref ToMail ## The email address used as the "to address", configured using parameters
        FROMADDRESS: !Ref FromMail ## The email address used as the "from address", configured using parameters
        APITOKEN: !Ref ClubhouseToken ## The API token for Clubhouse
        DAYS: 7 ## The number of days to search for upcoming stories
        OWNER: retgits ## The name of the person to search stories for
        WAVEFRONT_ENABLED: true ## Send metrics to Wavefront or not
        WAVEFRONT_URL: !Ref WavefrontURL ## The URL to connect to Wavefront
        WAVEFRONT_API_TOKEN: !Ref WavefrontToken ## The API token for Wavefront
```

The environment variables that have a `!Ref` in front of it, are configured using CloudFormation parameters. These parameters are at the top of the [template.yaml](./template.yaml) and are set in the Make targets.

## Contributing

[Pull requests](https://github.com/retgits/clubhouse-duedates/pulls) are welcome. For major changes, please open [an issue](https://github.com/retgits/clubhouse-duedates/issues) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

See the [LICENSE](./LICENSE) file in the repository