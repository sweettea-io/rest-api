package cloud

var CurrentCloud Cloud

func InitCloud(provider string, config map[string]string) {
  var err error

  switch provider {
  case AWS:
    CurrentCloud, err = NewAWSCloud(
      config["region"],
      config["sslCert"],
    )
  }

  if err != nil {
    panic(err)
  }
}