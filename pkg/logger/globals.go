package logger

import "github.com/sirupsen/logrus"

// CtxKeys is used to set global key->values that the logger need
var CtxKeys = map[string]string{}

// Logger main logger instance
var Logger *logrus.Logger
