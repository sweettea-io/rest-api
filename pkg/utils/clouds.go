package utils

type cloudsType struct {
  AWS string
}

type cloudsMapType map[string]bool

var Clouds = cloudsType{
  AWS: "aws",
}

var CloudsMap = cloudsMapType{
  Clouds.AWS: true,
}

// Check if a cloud is supported.
func IsValidCloud(cloud string) bool {
  return CloudsMap[cloud] == true
}