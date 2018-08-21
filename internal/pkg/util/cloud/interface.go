package cloud

type Cloud interface {
  Init() error
  SSLServiceLabels() map[string]string
}