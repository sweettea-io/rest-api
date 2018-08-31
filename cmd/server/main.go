package main

import (
  "fmt"
  "net/http"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/app/route"
  "github.com/sweettea-io/rest-api/internal/pkg/config"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
  "github.com/sweettea-io/rest-api/internal/pkg/util/dns"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Init cloud instance.
  cloud.InitCloud(app.Config.CloudProvider, map[string]string{
    "region": app.Config.AWSRegionName,
    "sslCert": app.Config.WildcardSSLCertArn,
  })

  // Init DNS service provider.
  dns.InitDNS(app.Config.DNSService, map[string]string{
    "hostedZoneId": app.Config.HostedZoneId,
  })

  // Build API routes.
  route.InitRouter(app.Config)

  // Format address to listen on.
  addr := fmt.Sprintf(":%v", app.Config.ServerPort)

  app.Log.Infof("Listening on port %v...\n", app.Config.ServerPort)

  // Start server.
  panic(http.ListenAndServe(addr, route.Router.GetRouter()))
}