package cloud

// Supported clouds.
const (
  AWS = "aws"
)

var validClouds = map[string]bool{
  AWS: true,
}

// IsValidCloud returns whether the provided cloud is a supported cloud.
func IsValidCloud(cloud string) bool {
  return validClouds[cloud] == true
}