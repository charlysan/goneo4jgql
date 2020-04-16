// Package logger provides logging functionality using Apache Format
// The following levels are supported:
//    "debug"
//    "info"
//    "warning"
//    "error"
//    "fatal"
//
// Levels should be specified using env var LOGGER_LEVEL
// Then the logger will only log entries with that severity or anything above it.
//   e.g.  LOGGER_LEVEL=INFO --> Will log anything that is info or above
//   (warn, error, fatal)
package logger

import (
	"os"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// LogLevel type used for specifying log level
type LogLevel string

// LogFormatter type used for specifying formatter
type LogFormatter string

// LogFields type used for specifying log fields
type LogFields map[string]interface{}

const (
	LevelDebug   LogLevel = "debug"
	LevelInfo    LogLevel = "info"
	LevelWarning LogLevel = "warning"
	LevelError   LogLevel = "error"
	LevelFatal   LogLevel = "fatal"
)

const (
	// FormatConsole formats logs output to a "human friendly format"
	FormatConsole LogFormatter = "console"
	// FormatLogstash formats logs output to logstash format (JSON)
	FormatLogstash LogFormatter = "logstash"
)

// InitializeLogger initializes a new logger with format and log level
func InitializeLogger() {
	logger := logrus.New()

	// TODO: add possibility to output to file in the future
	logger.SetOutput(os.Stderr)

	formatter := viper.GetString("LOGGER_FORMATTER")
	switch strings.ToLower(formatter) {
	case string(FormatConsole):
		logger.SetFormatter(&prefixed.TextFormatter{})
	case string(FormatLogstash):
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	level := viper.GetString("LOGGER_LEVEL")
	switch strings.ToLower(level) {
	case string(LevelDebug):
		logger.SetLevel(logrus.DebugLevel)
	case string(LevelInfo):
		logger.SetLevel(logrus.InfoLevel)
	case string(LevelWarning):
		logger.SetLevel(logrus.WarnLevel)
	case string(LevelError):
		logger.SetLevel(logrus.ErrorLevel)
	case string(LevelFatal):
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
	}

	Logger = logger
}

// parseArguments parses arguments
func parseArguments(args []interface{}) (string, map[string]interface{}) {
	msg := ""
	fields := map[string]interface{}{}
	fields["prefix"] = "main"

	for _, arg := range args {
		switch t := arg.(type) {
		case LogFields:
			for k, v := range t {
				fields[k] = v
			}
		case error:
			trace := errors.Wrap(arg, 0).ErrorStack()
			fields["error_msg"] = t.Error()
			fields["trace"] = trace
		case string:

			msg = t
		}

	}

	return msg, fields
}

// NewLogger returns a logger (Entry) that can be used by third party middlewares
func NewLogger() *logrus.Entry {
	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	_, fields := parseArguments([]interface{}{})
	return Logger.WithFields(fields)
}

// Debug logs a debug message with "debug" level
// It supports the following arguments:
//  - string (for main log message)
//  - error (for logging error trace)
//  - logger.LogFields (for custom fields besides the log message - e.g. item_uuid)
func Debug(args ...interface{}) {
	level := strings.ToLower(viper.GetString("LOGGER_LEVEL"))
	if level != string(LevelDebug) {
		return
	}

	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	msg, fields := parseArguments(args)
	Logger.WithFields(fields).Debug(msg)
}

// Info logs a debug message with "info" level
// It supports the following arguments:
//  - string (for main log message)
//  - error (for logging error trace)
//  - logger.LogFields (for custom fields besides the log message - e.g. item_uuid)
func Info(args ...interface{}) {
	level := strings.ToLower(viper.GetString("LOGGER_LEVEL"))
	if level != string(LevelInfo) && level != string(LevelDebug) {
		return
	}

	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	msg, fields := parseArguments(args)
	Logger.WithFields(fields).Info(msg)
}

// Warning logs a debug message with "Warning" level
// It supports the following arguments:
//  - string (for main log message)
//  - error (for logging error trace)
//  - logger.LogFields (for custom fields besides the log message - e.g. item_uuid)
func Warning(args ...interface{}) {
	level := strings.ToLower(viper.GetString("LOGGER_LEVEL"))
	if level == string(LevelError) || level == string(LevelFatal) {
		return
	}

	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	msg, fields := parseArguments(args)
	Logger.WithFields(fields).Warning(msg)
}

// Error logs a debug message with "Error" level
// It supports the following arguments:
//  - string (for main log message)
//  - error (for logging error trace)
//  - logger.LogFields (for custom fields besides the log message - e.g. item_uuid)
func Error(args ...interface{}) {
	level := strings.ToLower(viper.GetString("LOGGER_LEVEL"))
	if level == string(LevelFatal) {
		return
	}

	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	msg, fields := parseArguments(args)
	Logger.WithFields(fields).Error(msg)
}

// Fatal logs a debug message with "Error" level
// It supports the following arguments:
//  - string (for main log message)
//  - error (for logging error trace)
//  - logger.LogFields (for custom fields besides the log message - e.g. item_uuid)
func Fatal(args ...interface{}) {
	// Initialize logger if it hasn't been initialized yet
	if Logger == nil {
		InitializeLogger()
	}

	msg, fields := parseArguments(args)
	Logger.WithFields(fields).Fatal(msg)
}
