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
* [A Pulumi account](https://pulumi.com)

## Installation

### Get sources

To install you can clone this repository, and run `go get ./...` to get the [Go modules](./go.mod) this app relies on.

### Using Pulumi

Pulumi enables developers to write code in their favorite language, such as Go. This enables modern approaches to cloud applications and infrastructure without needing to learn yet-another YAML or DSL dialect.

### Email policy

To be able to send emails you'll need to allow AWS Lambda to access SES. You can follow the "_Least Privilege Configuration_" and create a new policy with the below content and the ARN to the list of policies.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "ses:GetIdentityVerificationAttributes",
                "ses:SendEmail",
                "ses:SendRawEmail",
                "ses:VerifyEmailIdentity"
            ],
            "Resource": "arn:aws:ses:us-west-2:<ACCOUNTID>:identity/<EMAILADDRESS>",
            "Effect": "Allow"
        }
    ]
}
```

The other option is to add "**arn:aws:iam::aws:policy/AmazonSESFullAccess**", which is an AWS managed policy with a lot more privileges but saves you from creating a new policy.

## Configuration

Inside the [Pulumi.dev.yaml](./pulumi/Pulumi.dev.yaml), there are a few configuration options that you can set:

```yaml
config:
  aws:profile: default ## The AWS CLI profile you want to use
  aws:region: us-west-2 ## The AWS region you want to deploy to
  clubhouse-duedates:config:
    s3bucket: <my-bucket> ## the S3 bucket your code will be uploaded to
    policies:
      - <ARN of Send Email policy>
      - arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
      - arn:aws:iam::aws:policy/service-role/AWSLambdaRole
    envvars:
      - REGION/us-west-2
      - TOADDRESS/<your configured SES email addres>
      - FROMADDRESS/<your configured SES email addres>
      - APITOKEN/<your clubhouse API token>
      - DAYS/7 ## The number of days to look ahead for stories
      - OWNER/<your name> ## Your name
      - WAVEFRONT_ENABLED/true
      - WAVEFRONT_URL/WavefrontURL ## The Wavefront URL of your Wavefront instance
      - WAVEFRONT_API_TOKEN/WavefrontToken ## The Wavefront token to connect to Wavefront
    tags:
      author: <your name>
      feature: clubhouse-duedates
      region: us-west-2
      team: <your team>
      version: v1
```

## Contributing

[Pull requests](https://github.com/retgits/clubhouse-duedates/pulls) are welcome. For major changes, please open [an issue](https://github.com/retgits/clubhouse-duedates/issues) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

See the [LICENSE](./LICENSE) file in the repository
