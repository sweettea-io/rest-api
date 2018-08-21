package cloud

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
)

type AWSCloud struct {
  Client  *session.Session
  Region  string
  SSLCert string
}

func NewAWSCloud(region string, sslCert string) (*AWSCloud, error) {
  // Create a new AWS Session.
  sess, err := session.NewSession(&aws.Config{
    Region: &region,
    Credentials: credentials.NewEnvCredentials(),
  })

  if err != nil {
    return nil, fmt.Errorf("error initializing new AWS session: %s", err.Error())
  }

  awsCloud := &AWSCloud{
    Client: sess,
    Region: region,
    SSLCert: sslCert,
  }

  return awsCloud, nil
}

func (a *AWSCloud) GetClient() *session.Session {
  return a.Client
}

func (a *AWSCloud) SSLServiceLabels() map[string]string {
  return map[string]string{
    "service.beta.kubernetes.io/aws-load-balancer-ssl-cert": a.SSLCert,
    "service.beta.kubernetes.io/aws-load-balancer-ssl-ports": "443",
  }
}