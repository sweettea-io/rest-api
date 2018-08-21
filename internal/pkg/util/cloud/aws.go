package cloud

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/sweettea-io/rest-api/internal/app"
)

type AWSCloud struct {
  Client *session.Session
}

func (a *AWSCloud) Init() error {
  // Create a new AWS Session.
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String(app.Config.AWSRegionName),
    Credentials: credentials.NewEnvCredentials(),
  })

  if err != nil {
    return fmt.Errorf("error initializing new AWS session: %s", err.Error())
  }

  // Store session object as Client.
  a.Client = sess

  return nil
}

func (a *AWSCloud) SSLServiceLabels() map[string]string {
  return map[string]string{
    "service.beta.kubernetes.io/aws-load-balancer-ssl-cert": app.Config.WildcardSSLCertArn,
    "service.beta.kubernetes.io/aws-load-balancer-ssl-ports": "443",
  }
}