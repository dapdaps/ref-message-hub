package log

import (
	"github.com/sirupsen/logrus"
	"ref-message-hub/common"
	"ref-message-hub/common/referror"
	"strings"
	"time"
)

var (
	xlog = logrus.New()
)

type Config struct {
	Dir   string
	Level string
	//Alarm ecode.Alarm
}

func Init(conf *Config, debug bool) {
	xlog.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	if conf.Dir != "" {
		pathMap := PathMap{
			logrus.InfoLevel:  conf.Dir + "/info.log",
			logrus.ErrorLevel: conf.Dir + "/error.log",
		}
		xlog.Hooks.Add(NewLocalHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	}
	if debug {
		xlog.SetLevel(logrus.DebugLevel)
	} else {
		if conf.Level != "" {
			slevel := strings.ToLower(conf.Level)
			switch slevel {
			case "fatal":
				xlog.SetLevel(logrus.FatalLevel)
			case "error":
				xlog.SetLevel(logrus.ErrorLevel)
			case "warn":
				xlog.SetLevel(logrus.WarnLevel)
			case "info":
				fallthrough
			default:
				xlog.SetLevel(logrus.InfoLevel)
			}
		} else {
			xlog.SetLevel(logrus.InfoLevel)
		}
	}
	//xlog.Hooks.Add(NewPrometheusHook(conf.Alarm))
}

func Debug(format string, args ...interface{}) {
	xlog.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	xlog.Infof(format, args...)
}

func InfoUUID(uuid string, format string, args ...interface{}) {
	xlog.WithField("server_uuid", uuid).Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	var entry *logrus.Entry
	entry, args = EntryWithMetric(args...)
	entry.Warnf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	var entry *logrus.Entry
	entry, args = EntryWithMetric(args...)
	entry.Fatalf(format, args...)
}

func LevelEnabled(level logrus.Level) bool {
	return xlog.IsLevelEnabled(level)
}

func Entry() *logrus.Entry {
	return logrus.NewEntry(xlog)
}

func Error(format string, args ...interface{}) {
	var entry *logrus.Entry
	entry, args = EntryWithMetric(args...)
	entry.Errorf(format, args...)
}

func EntryWithMetric(args ...interface{}) (*logrus.Entry, []interface{}) {
	var ecode int
	var alertText string
	for i, arg := range args {
		if err, ok := arg.(*referror.Error); ok {
			ecode = err.Code
			alertText = err.Message
			args[i] = alertText
			break
		}
	}

	// requestId, b := ctx.Value(RequestID).(string)
	// if !b {
	// 	requestId = "no-requestId"
	// }

	requestId := ""

	return xlog.WithField(Ecode, ecode).
		WithField(RequestID, requestId).
		WithField(AlertText, alertText).
		WithField("stack", common.Stack(2, 1000)), args
}
