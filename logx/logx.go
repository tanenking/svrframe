package logx

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/constants"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// level ErrorLevel-DebugLevel
func InitLogx() {

	if inited {
		return
	}
	inited = true

	level := DebugLevel
	if !constants.IsDebug() {
		log_level := os.Getenv(constants.Env_LogLevel)
		if len(log_level) > 0 {
			l, e := strconv.Atoi(log_level)
			if e == nil {
				level = logrus.Level(l)
			}
		}
		if level < ErrorLevel || level >= MaxLevel {
			level = InfoLevel
		}
	}

	fullpath = getLogPath()

	fullpath += "/" + constants.Service_Type

	filename = fmt.Sprintf("log-%d.info", config.GetServiceInfo().ServiceID)
	file_info := fullpath + "/" + filename

	writer_info, _ := rotatelogs.New(
		file_info+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(file_info),
		rotatelogs.WithMaxAge(file_WithMaxAge),
		rotatelogs.WithRotationCount(file_WithRotationCount),
		rotatelogs.WithRotationTime(file_WithRotationTime),
		rotatelogs.WithRotationSize(file_WithRotationSize),
	)
	mw_info := io.MultiWriter(writer_info)
	//if constants.IsDebug() || isWindowsSystem() {
	//mw_info = io.MultiWriter(mw_info, os.Stdout)
	//}

	logInfo.SetOutput(mw_info)

	logInfo.SetLevel(level)

	formatter := &formatter{
		//ForceColors:               true,
		//EnvironmentOverrideColors: true,
		//FullTimestamp:   true,
		//TimestampFormat: constants.TimeFormatString,
		//DisableSorting:            true,
		//DisableLevelTruncation:    true,
		//PadLevelText: true,
	}

	logInfo.SetFormatter(formatter)

	go update()
}

func GetTraceFile() string {
	return fullpath + "/trace"
}
func GetLogFullPath() string {
	return fullpath
}
func GetLogFileName() string {
	return filename
}

func GetLoggerWriter() io.Writer {
	return logInfo.Out
}

func TraceBack() {
	if !inited {
		logrus.Errorf(string(debug.Stack()))
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   string(debug.Stack()),
		ctx:   ctx,
		level: DebugLevel,
	}
	buffer <- t
}

func Debugln(args ...interface{}) {
	if !inited {
		logrus.Debugln(args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintln(args...),
		ctx:   ctx,
		level: DebugLevel,
	}
	buffer <- t
}
func DebugF(msg string, args ...interface{}) {
	if !inited {
		logrus.Debugf(msg, args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintf(msg, args...),
		ctx:   ctx,
		level: DebugLevel,
	}
	buffer <- t
}

func Infoln(args ...interface{}) {
	if !inited {
		logrus.Infoln(args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintln(args...),
		ctx:   ctx,
		level: InfoLevel,
	}
	buffer <- t
}
func InfoF(msg string, args ...interface{}) {
	if !inited {
		logrus.Infof(msg, args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintf(msg, args...),
		ctx:   ctx,
		level: InfoLevel,
	}
	buffer <- t
}

func Warnln(args ...interface{}) {
	if !inited {
		logrus.Warnln(args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintln(args...),
		ctx:   ctx,
		level: WarnLevel,
	}
	buffer <- t
}
func WarnF(msg string, args ...interface{}) {
	if !inited {
		logrus.Warnf(msg, args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintf(msg, args...),
		ctx:   ctx,
		level: WarnLevel,
	}
	buffer <- t
}

func Errorln(args ...interface{}) {
	if !inited {
		logrus.Errorln(args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintln(args...),
		ctx:   ctx,
		level: ErrorLevel,
	}
	buffer <- t
}
func ErrorF(msg string, args ...interface{}) {
	if !inited {
		logrus.Errorf(msg, args...)
		return
	}
	ctx := getContext(3)
	t := &log_t{
		msg:   fmt.Sprintf(msg, args...),
		ctx:   ctx,
		level: ErrorLevel,
	}
	buffer <- t
}
