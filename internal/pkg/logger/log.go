package logger

import "github.com/Sirupsen/logrus"

func NewLogger() *logrus.Logger {
  // Add whatever custom log formatting to the logger here.
  return logrus.New()
}