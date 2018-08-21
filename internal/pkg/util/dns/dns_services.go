package dns

// Supported DNS Services.
const (
  Route53 = "route53"
)

var validDNSServices = map[string]bool{
  Route53: true,
}

// IsValidDNS returns whether the provided DNS service is supported.
func IsValidDNS(dnsService string) bool {
  return validDNSServices[dnsService] == true
}