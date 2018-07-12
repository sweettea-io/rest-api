package utils

type cloudsType struct {
  AWS string
}

var Clouds = cloudsType{
  AWS: "aws",
}

var validClouds = map[string]bool {
  Clouds.AWS: true,
}

// Check if a cloud is supported.
func IsValidCloud(cloud string) bool {
  return validClouds[cloud] == true
}