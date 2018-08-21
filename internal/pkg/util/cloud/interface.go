package cloud

type Cloud interface {
  SSLServiceLabels() map[string]string
}