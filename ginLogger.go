package goginlogger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ILogrus interface {
	Debug(*gin.Context, ...interface{})
	Info(*gin.Context, ...interface{})
	Warning(*gin.Context, ...interface{})
	Error(*gin.Context, ...interface{})
}

type logrus struct {
	logger        map[string]*log.Logger
	contextFields []ContextField
}

var logLevels = map[string]log.Level{
	"Info":  log.InfoLevel,
	"Debug": log.DebugLevel,
	"Warn":  log.WarnLevel,
	"Error": log.ErrorLevel,
}

type LoggerConfig struct {
	Env    string
	LogDir string
}

// This is an alias of type string to use as context parameter
// Example Usage: this is the example of gin middleware to set ContextField fields
//
//	func(ctx *gin.Context) {
//			newId := uuid.New().String()
//			requestId := ctx.Request.Header.Get("X-Request-Id")
//			if requestId == "" {
//				ctx.Set("X-Request-Id", newId)
//			}
//			ctx.Next()
//		}
type ContextField string

// NewGinLogger is higher order function
// This function is designed to create and return a logger function tailored specifically for use within a Gin web framework environment.
// It create four log files in given dir eg. Info.log, Debug.log ...
// NewGinLogger returns a new (func() ILogrus).
// It requires the following parameters:
// 1. config (*dig.LoggerConfig): to decide the logger folder and its ev.
// 2. contextFields (...ContextField): receives n number of parameter of type string to log

// Example Usage:
//
//	logger := NewGinLogger(logConfig, ContextField("X-Request-Id"), ContextField("X-User-Id"))(),
//	// Use the 'logger' instance to use different logging method Info, Debug, Error, Warning.

func NewGinLogger(config LoggerConfig, contextFields ...ContextField) func() ILogrus {

	return func() ILogrus {
		logInstance := &logrus{
			logger:        make(map[string]*log.Logger),
			contextFields: contextFields,
		}

		if _, err := os.Stat(config.LogDir); os.IsNotExist(err) {
			// log dir does not exist
			if err := os.Mkdir(config.LogDir, 0755); err != nil {
				log.Fatal("Error occured while creating logger dir", err)
			}
		}

		for severity, logLevel := range logLevels {
			lg := log.New()
			lg.SetFormatter(&log.JSONFormatter{})
			lg.SetLevel(logLevel)
			file, err := os.OpenFile(
				fmt.Sprintf("%s/%s.log", config.LogDir, severity), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666,
			)
			if err == nil {
				if config.Env == "PROD" || config.Env == "STAGE" {
					lg.SetOutput(file)
				} else {
					lg.SetOutput(io.MultiWriter(file, os.Stdout))
				}

			} else {
				log.Fatal("Error occured while creating logger file", err)
			}
			logInstance.logger[severity] = lg
		}

		return logInstance
	}

}

func (l *logrus) Debug(ctx *gin.Context, args ...interface{}) {
	if ctx != nil {
		l.logger["Debug"].WithFields(getLogFields(ctx, l.contextFields)).Info(args...)
	} else {
		l.logger["Debug"].Info(args...)
	}
}

func (l *logrus) Info(ctx *gin.Context, args ...interface{}) {
	if ctx != nil {
		l.logger["Info"].WithFields(getLogFields(ctx, l.contextFields)).Info(args...)
	} else {
		l.logger["Info"].Info(args...)
	}

}
func (l *logrus) Warning(ctx *gin.Context, args ...interface{}) {
	if ctx != nil {
		l.logger["Warn"].WithFields(getLogFields(ctx, l.contextFields)).Warn(args...)
	} else {
		l.logger["Warn"].Warn(args...)
	}

}
func (l *logrus) Error(ctx *gin.Context, args ...interface{}) {
	if ctx != nil {
		l.logger["Error"].WithFields(getLogFields(ctx, l.contextFields)).Warn(args...)
	} else {
		l.logger["Error"].Warn(args...)
	}
}

// this method is used to fetch all the fields which are provided while logger instance
func getLogFields(ctx *gin.Context, contextFields []ContextField) log.Fields {
	fields := log.Fields{}
	for _, contextField := range contextFields {
		if val, found := ctx.Get((string)(contextField)); found {
			fields[strings.ToLower((string)(contextField))] = val
		}

	}
	return fields
}
