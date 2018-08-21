package dns

type DNS interface {
  UpsertRR(rrType string, name string, values []string, ttl uint) error
}