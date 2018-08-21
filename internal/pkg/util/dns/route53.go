package dns

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/aws/aws-sdk-go/service/route53"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
)

type Route53DNS struct {
  Client *route53.Route53
  HostedZoneID string
}

func NewRoute53() *Route53DNS {
  // Route53 requires the AWS cloud provider to be configured.
  currCloud, ok := cloud.CurrentCloud.(*cloud.AWSCloud)

  if !ok {
    panic(fmt.Errorf("Can't configure Route53 as DNS Service without AWS as cloud provider..."))
  }

  if app.Config.HostedZoneId == "" {
    panic(fmt.Errorf("HostedZoneId config can't be blank when using Route53 as DNS Service"))
  }

  return &Route53DNS{
    Client: route53.New(currCloud.Client),
    HostedZoneID: app.Config.HostedZoneId,
  }
}

func (r *Route53DNS) UpsertRR(rrType string, name string, values []string, ttl int64) error {
  // Create var from const of upsert action (need to take pointer)
  action := route53.ChangeActionUpsert

  // Create RR values from string values.
  var rrValues []*route53.ResourceRecord
  for _, val := range values {
    rrValues = append(rrValues, &route53.ResourceRecord{Value: &val})
  }

  // Prep changes for Route53 API call.
  changes := &route53.ChangeResourceRecordSetsInput{
    HostedZoneId: &r.HostedZoneID,
    ChangeBatch: &route53.ChangeBatch{
      Changes: []*route53.Change{{
        Action: &action,
        ResourceRecordSet: &route53.ResourceRecordSet{
          Name: &name,
          Type: &rrType,
          TTL: &ttl,
          ResourceRecords: rrValues,
        },
      }},
    },
  }

  if _, err := r.Client.ChangeResourceRecordSets(changes); err != nil {
    return fmt.Errorf("Error upserting RRs on Route53: %s", err.Error())
  }

  return nil
}