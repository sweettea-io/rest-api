package dns

import "github.com/sweettea-io/rest-api/internal/app"

var CurrentDNS DNS

func InitDNS() {
  switch app.Config.DNSService {
  case Route53:
    CurrentDNS = NewRoute53()
  }
}