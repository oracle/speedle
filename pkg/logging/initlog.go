//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package logging

import (
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/oracle/speedle/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	FORMATTER_TEXT = "text"
	FORMATTER_JSON = "json"
)

// LogConfig is the struct for the logger
type LogConfig struct {
	Level           string             `json:"level,omitempty"`
	Formatter       string             `json:"formatter,omitempty"`
	SetReportCaller bool               `json:"setReportCaller,omitempty"`
	RotationConfig  *lumberjack.Logger `json:"rotationConfig,omitempty"`
}

// InitLog initializes the standard logger instance
func InitLog(cfg *LogConfig) error {
	// Configure the log level, defaults to "WARN"
	if cfg != nil && len(cfg.Level) != 0 {
		logLevel, err := log.ParseLevel(cfg.Level)
		if err != nil {
			return errors.Wrapf(err, errors.LoggingError, "failed to parse log level: %q", cfg.Level)
		}
		log.SetLevel(logLevel)
	}

	// Configure the reportCaller, defaults to false
	if cfg != nil && cfg.SetReportCaller {
		log.SetReportCaller(cfg.SetReportCaller)
	}

	// Configure the log formatter, defaults to ASCII formatter
	if cfg != nil && len(cfg.Formatter) != 0 {
		switch strings.ToLower(cfg.Formatter) {
		case FORMATTER_TEXT:
			log.SetFormatter(&log.TextFormatter{})
		case FORMATTER_JSON:
			log.SetFormatter(&log.JSONFormatter{})
		default:
			return errors.Errorf(errors.LoggingError, "unknown formatter type: %q", cfg.Formatter)
		}
	}

	// Configure the file roration PluginConfig, defaults to os.Stderr
	if cfg != nil && cfg.RotationConfig != nil {
		log.SetOutput(cfg.RotationConfig)
	}

	return nil
}

// InitLogInstance initializes the specific logger instance
func InitLogInstance(logger *log.Logger, cfg *LogConfig) error {
	if logger == nil {
		return errors.New(errors.LoggingError, "the logger instance is nil")
	}

	// Configure the log level, defaults to "WARN"
	if cfg != nil && len(cfg.Level) != 0 {
		logLevel, err := log.ParseLevel(cfg.Level)
		if err != nil {
			return errors.Wrapf(err, errors.LoggingError, "failed to parse log level: %q", cfg.Level)
		}

		logger.SetLevel(logLevel)
	}

	// Configure the reportCaller, defaults to false
	if cfg != nil && cfg.SetReportCaller {
		logger.SetReportCaller(cfg.SetReportCaller)
	}

	// Configure the log formatter, defaults to ASCII formatter
	if cfg != nil && len(cfg.Formatter) != 0 {
		switch strings.ToLower(cfg.Formatter) {
		case FORMATTER_TEXT:
			logger.Formatter = &log.TextFormatter{}
		case FORMATTER_JSON:
			logger.Formatter = &log.JSONFormatter{}
		default:
			return errors.Errorf(errors.LoggingError, "unknown formatter type: %q", cfg.Formatter)
		}
	}

	// Configure the file roration PluginConfig, defaults to os.Stderr
	if cfg != nil && cfg.RotationConfig != nil {
		logger.Out = cfg.RotationConfig
	}

	return nil
}
