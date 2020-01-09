package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

const (
	shell     = "sh"
	shellFlag = "-c"
)

// Config contains the toplevel structure of the configuration variables
type Config struct {
	S3Bucket string
	Policies []string
	EnvVars  []string
	Tags     Tags
}

// Tags are used for each individual resource so they can be found using the Resource Groups service in the AWS Console
type Tags struct {
	Author  string
	Feature string
	Region  string
	Team    string
	Version string
}

// ToMap marshals the Tags into a map[string]string as required by Pulumi
func (t Tags) ToMap() map[string]string {
	tags := make(map[string]string)

	tags["version"] = t.Version
	tags["author"] = t.Author
	tags["team"] = t.Team
	tags["feature"] = t.Feature
	tags["region"] = t.Region

	return tags
}

var cfg Config

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new config bag with the given context and empty namespace and loads
		// configuration values into a Config object, or panics if unable to do so.
		config.New(ctx, "").RequireObject("config", &cfg)

		// Build the executable for AWS Lambda
		if err := Run("GOOS=linux GOARCH=amd64 go build -o ./bin/getstories ./getstories"); err != nil {
			return fmt.Errorf("Error building code: %s", err.Error())
		}

		// Create a zipfile
		if err := Run("zip -r -j ./bin/getstories.zip ./bin/getstories"); err != nil {
			return fmt.Errorf("Error creating zipfile: %s", err.Error())
		}

		// Upload the zipfile to S3
		if err := Run(fmt.Sprintf("aws s3 cp ./bin/getstories.zip s3://%s/getstories.zip", cfg.S3Bucket)); err != nil {
			return fmt.Errorf("Error uploading zipfile: %s", err.Error())
		}

		// Create the IAM Role for the Lambda function
		roleArgs := &iam.RoleArgs{
			AssumeRolePolicy:    `{"Version": "2012-10-17","Statement": [{"Action": "sts:AssumeRole","Principal": {"Service": "lambda.amazonaws.com"},"Effect": "Allow","Sid": ""}]}`,
			Description:         "IAM Role for Clubhouse Duedates",
			ForceDetachPolicies: false,
			MaxSessionDuration:  3600,
			Name:                "ClubhouseDuedatesRole",
			Tags:                cfg.Tags.ToMap(),
		}

		role, err := iam.NewRole(ctx, "ClubhouseDuedatesRole", roleArgs)
		if err != nil {
			return err
		}

		// Attach all policies to the new role
		for idx, policy := range cfg.Policies {
			rolePolicyAttachmentArgs := &iam.RolePolicyAttachmentArgs{
				PolicyArn: policy,
				Role:      role.ID(),
			}

			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("Policy Attachment %d", idx), rolePolicyAttachmentArgs)
			if err != nil {
				return err
			}
		}

		// Create a map of environment variables
		variables := make(map[string]interface{})
		for _, envvar := range cfg.EnvVars {
			parts := strings.Split(envvar, "/")
			variables[parts[0]] = parts[1]
		}
		environment := make(map[string]interface{})
		environment["variables"] = variables

		// Create the AWS Lambda function
		// The set of arguments for constructing a Function resource.
		functionArgs := &lambda.FunctionArgs{
			Description: "Get stories that near their due date from Clubhouse and send them to my email",
			Runtime:     "go1.x",
			Name:        "GetStories",
			MemorySize:  256,
			Timeout:     180,
			Handler:     "getstories",
			Environment: environment,
			S3Bucket:    cfg.S3Bucket,
			S3Key:       "getstories.zip",
			Role:        role.Arn(),
		}

		// NewFunction registers a new resource with the given unique name, arguments, and options.
		function, err := lambda.NewFunction(ctx, "GetStoriesFunction", functionArgs)
		if err != nil {
			return err
		}

		// Export the function ARN as an output of the Pulumi stack
		ctx.Export("GetStoriesFunctionARN", function.Arn())

		eventRuleArgs := &cloudwatch.EventRuleArgs{
			Description:        "Trigger for Clubhouse Duedates - GetStories",
			Name:               "GetStoriesTrigger",
			IsEnabled:          true,
			Tags:               cfg.Tags.ToMap(),
			ScheduleExpression: "cron(0 13 ? * * *)",
		}

		rule, err := cloudwatch.NewEventRule(ctx, "GetStoriesEventRule", eventRuleArgs)
		if err != nil {
			return err
		}

		eventTargetArgs := &cloudwatch.EventTargetArgs{
			Arn:  function.Arn(),
			Rule: rule.ID(),
		}

		_, err = cloudwatch.NewEventTarget(ctx, "GetStoriesEventTarget", eventTargetArgs)
		if err != nil {
			return err
		}

		return nil
	})
}

// Run starts the specified command and waits for it to complete.
func Run(args string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(shell, shellFlag, args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path.Join(dir, "..")
	return cmd.Run()
}
