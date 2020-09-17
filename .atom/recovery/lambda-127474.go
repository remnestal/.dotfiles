package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/slack-go/slack"
)

const (
	ENV_SLACK_ACCESS_TOKEN string = "slack_access_token"
	ENV_SLACK_CHANNEL      string = "slack_channel"
	ENV_RSA_BIT_SIZE       string = "rsa_bit_size"
)

// ErrOperationNotImplemented is the default error to return when the operation
// specified by the Secrets Manager can't be interpreted.
var ErrOperationNotImplemented error = errors.New("operation not implemented")

// Event defines the schema of the parameters sent by the Secrets Manager.
type Event struct {
	Step               string `json:"Step"`
	SecretId           string `json:"SecretId"`
	ClientRequestToken string `json:"ClientRequestToken"`
}

// slack_report is an abstraction for writing a message to Slack using the SDK.
func slack_report(secret_name string, message string) error {
	// Grab the slack bot access-token handle from environment.
	token, exists := os.LookupEnv(ENV_SLACK_ACCESS_TOKEN)
	if !exists {
		return fmt.Errorf("environment variable %q not set", ENV_SLACK_ACCESS_TOKEN)
	}
	// Grab the slack channel handle from environment.
	channel, exists := os.LookupEnv(ENV_SLACK_CHANNEL)
	if !exists {
		return fmt.Errorf("environment variable %q not set", ENV_SLACK_CHANNEL)
	}
	// Create a Slack client and post the message.
	api := slack.New(token)
	_, _, err := api.PostMessage(channel,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText(message, false),
	)
	return err
}

// report_success is an abstraction for reporting a successful key rotation.
func report_success(secret_name string) {
	if err := slack_report(secret_name, fmt.Sprintf("♻️ *Rotated* `%s`", secret_name)); err != nil {
		log.Printf("unable to post rotation success to Slack: %v", err)
	}
}

// report_failure is an abstraction for reporting a failed key rotation.
func report_failure(secret_name string) {
	if err := slack_report(secret_name, fmt.Sprintf("💥 *Rotation failed* `%s`", secret_name)); err != nil {
		log.Printf("unable to post rotation success to Slack: %v", err)
	}
}

func main() {
	lambda.Start(func(ctx context.Context, event Event) (string, error) {
		switch event.Step {
		case "createSecret":
			// createSecret is the operation invoked by the Secrets Manager when a
			// rotation is issued. This step involves creating a new key value and
			// pushing it as a new version of the secret.
			if err := rotate_rsa(ctx, event); err != nil {
				report_failure(event.SecretId)
				return err.Error(), err
			} else {
				report_success(event.SecretId)
				return "Rotation succeeded", nil
			}

		default:
			return fmt.Sprintf("unsupported: %s", event.Step), ErrOperationNotImplemented
		}
	})
}

// rotate_rsa generates a new RSA private key and pushes it to the specified secret
// entry in the Secrets Manager.
func rotate_rsa(ctx context.Context, event Event) error {
	// Create a new Secrets Manager client.
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("unable to start new AWS session: %w", err)
	}
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion("eu-west-1"))

	// Generate a new RSA key pair.
	bitsize_string, exists := os.LookupEnv(ENV_RSA_BIT_SIZE)
	if !exists {
		return fmt.Errorf("environment variable %q not set", ENV_RSA_BIT_SIZE)
	}
	bitsize, err := strconv.ParseInt(bitsize_string, 10, 32)
	if err != nil {
		return fmt.Errorf("environment variable %q could not parsed as an integer", ENV_RSA_BIT_SIZE)
	}
	key, err := rsa.GenerateKey(rand.Reader, int(bitsize))
	if err != nil {
		return fmt.Errorf("unable to generate new RSA key pair: %w", err)
	}

	// Roll out the key as new version in the Secrets Manager.
	_, err = svc.PutSecretValueWithContext(ctx, &secretsmanager.PutSecretValueInput{
		SecretBinary:       x509.MarshalPKCS1PrivateKey(key),
		SecretId:           aws.String(event.SecretId),
		ClientRequestToken: aws.String(event.ClientRequestToken),
	})
	if err != nil {
		return fmt.Errorf("unable to put new version of secret value: %w", err)
	}

	return nil
}
