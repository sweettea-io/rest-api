package dns

import "github.com/aws/aws-sdk-go/service/route53"

type Route53DNS struct {}

func (r *Route53DNS) UpsertRR(rrType string, name string, values []string, ttl uint) error {
  return nil
}