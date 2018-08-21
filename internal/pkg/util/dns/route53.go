package dns

type Route53DNS struct {}

func (r *Route53DNS) UpsertRR(rrType string, name string, values []string, ttl uint) error {
  return nil
}