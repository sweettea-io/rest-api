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

func NewAWSCloud() (*AWSCloud, error) {
  // Create a new AWS Session.
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String(app.Config.AWSRegionName),
    Credentials: credentials.NewEnvCredentials(),
  })

  if err != nil {
    return nil, fmt.Errorf("error initializing new AWS session: %s", err.Error())
  }

  return &AWSCloud{Client: sess}, nil
}

func (a *AWSCloud) GetClient() *session.Session {
  return a.Client
}

func (a *AWSCloud) SSLServiceLabels() map[string]string {
  return map[string]string{
    "service.beta.kubernetes.io/aws-load-balancer-ssl-cert": app.Config.WildcardSSLCertArn,
    "service.beta.kubernetes.io/aws-load-balancer-ssl-ports": "443",
  }
}