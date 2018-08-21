package dns

type DNS interface {
  UpsertRR() error
}