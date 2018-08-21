package dns

var CurrentDNS DNS

func InitDNS(dnsService string, config map[string]string) {
  switch dnsService {
  case Route53:
    CurrentDNS = NewRoute53(config["hostedZoneId"])
  }
}