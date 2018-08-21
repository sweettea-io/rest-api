package cloud

import "github.com/sweettea-io/rest-api/internal/app"

type AWSCloud struct {}

func (aws *AWSCloud) SSLServiceLabels() map[string]string {
  return map[string]string{
    "service.beta.kubernetes.io/aws-load-balancer-ssl-cert": app.Config.WildcardSSLCertArn,
    "service.beta.kubernetes.io/aws-load-balancer-ssl-ports": "443",
  }
}