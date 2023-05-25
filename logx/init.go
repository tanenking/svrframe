package logx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tanenking/svrframe/constants"

	"github.com/sirupsen/logrus"
)

type log_t struct {
	msg   string
	ctx   context.Context
	level logrus.Level
}

var (
	logInfo  *logrus.Logger
	fullpath string
	filename string
	inited   bool
	buffer   chan *log_t
)

var (
	//windows console color
	proc        interface{}
	closeHandle interface{}
)

var (
	pid int
)

const (
	ErrorLevel logrus.Level = iota + logrus.ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	MaxLevel
)
const (
	//文件保留2周
	file_WithMaxAge = time.Duration(time.Hour * 24 * 7 * 2)
	//文件一天一切换
	file_WithRotationTime = time.Duration(time.Hour * 24)
	//
	file_WithRotationCount = 0
	//
	file_WithRotationSize = -1
	//消息队列长度
	msg_queue_max_len = 4096
)

func init() {
	logInfo = logrus.New()
	fullpath = ""
	filename = "log.info"
	inited = false
	pid = os.Getpid()

	initKernel32()

	buffer = make(chan *log_t, msg_queue_max_len)
}

func isLogRuntime() bool {
	log_runtime := os.Getenv(constants.Env_LogRuntime)
	log_runtime = strings.ToLower(log_runtime)
	if len(log_runtime) > 0 && log_runtime == "true" {
		return true
	}
	return false
}

func getLogPath() string {
	var _exists bool = false
	path := os.Getenv(constants.Env_LogPath)
	if len(path) <= 0 {
		gpath := os.Getenv(constants.Env_Gopath)
		if constants.IsDebug() && len(gpath) > 0 {
			path = gpath + "/log/" + constants.ProjectName
		} else {
			p, _ := os.Getwd()
			path = p + "/log/" + constants.ProjectName
		}
	} else {
		if !filepath.IsAbs(path) {
			// fmt.Printf("env log_path must need abs path")
			// panic("env log_path must need abs path")
			p, _ := os.Getwd()
			path = p + "/" + path
		}
		var is bool
		is, _exists = isDir(path)
		if _exists && !is {
			fmt.Printf("env log_path must need a dir")
			panic("env log_path must need a dir")
		}
	}
	if !_exists {
		//创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	return path
}

// 判断所给路径是否为文件夹
func isDir(path string) (is bool, exists bool) {
	is = false
	exists = false

	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		exists = true
		return
	}
	exists = true
	is = s.IsDir()
	return
}

func getRuntime(step int) string {

	if !isLogRuntime() {
		return ""
	}
	//function := "?"
	file := "?"
	line := 0

	_, file, line, ok := runtime.Caller(step)
	if !ok {
		//return fmt.Sprintf("%s %s:%d", file, function, line)
		return fmt.Sprintf("%s:%d", file, line)
	}
	//f := runtime.FuncForPC(pc)
	//if f != nil {
	//	function = f.Name()
	//}

	fName := filepath.Base(file)
	//return fmt.Sprintf("%s %s:%d", fName, function, line)
	return fmt.Sprintf("%s:%d", fName, line)
}

func getContext(_call_step int) context.Context {
	ctx := context.WithValue(context.Background(), "runtime", getRuntime(_call_step))
	//ctx = context.WithValue(ctx, "pid", pid)
	//ctx := context.WithValue(context.Background(), "pid", os.Getpid())

	return ctx
}

func update() {
	defer close(buffer)
	exit_ch := constants.GetServiceStopListener().AddListener()
	for {
		select {
		case <-exit_ch.Done():
			return
		case msg, ok := <-buffer:
			if !ok || msg == nil {
				logrus.Errorln("log buffer 关闭")
			} else {
				if msg.level == InfoLevel {
					logInfo.WithContext(msg.ctx).Infoln(msg.msg)
					colorPrint(msg.msg, green)
				} else if msg.level == WarnLevel {
					logInfo.WithContext(msg.ctx).Warnln(msg.msg)
					colorPrint(msg.msg, yellow)
				} else if msg.level == ErrorLevel {
					logInfo.WithContext(msg.ctx).Errorln(msg.msg)
					colorPrint(msg.msg, red)
				} else {
					logInfo.WithContext(msg.ctx).Debugln(msg.msg)
					colorPrint(msg.msg, gray)
				}
			}
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
