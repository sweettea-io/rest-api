package dns

type DNS interface {
  UpsertRR(rrType string, name string, values []string, ttl int64) error
}