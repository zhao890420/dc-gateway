package common

/**
 * @Author zhaoguang
 * @Description zap日志
 * @Date 4:31 下午 2021/2/21
 **/
import (
	"github.com/google/martian/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LogConfig struct {
	Logger Logger
}

// Logger 参数
type Logger struct {
	Mode                                 string
	Level                                string
	Name, Alone, Path, UUIDPath, ErrPath string
	MaxSize                              int
	MaxBackups                           int
	MaxAge                               int
	Compress                             bool
	//日志是否为json结构
	Json bool
}

// New 创建新日志客户端
func New(l *Logger) (*zap.Logger, error) {

	// 仅打印Error级别以上的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 获取日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(l.Level)); err != nil {
		return nil, err
	}
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	if l.ErrPath == "" {
		l.ErrPath = "error.log"
	}
	errHook := lumberjack.Logger{
		Filename:   l.ErrPath,
		MaxSize:    l.MaxSize,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		Compress:   l.Compress,
		LocalTime:  true,
	}

	hook := lumberjack.Logger{
		Filename:   l.Path,
		MaxSize:    l.MaxSize,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		Compress:   l.Compress,
		LocalTime:  true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据环境设置不同日志格式
	var encoder zapcore.Encoder
	if !l.Json {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(&hook), lowPriority),
		zapcore.NewCore(encoder, zapcore.AddSync(&errHook), highPriority),
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	return zap.New(core, caller), nil
}
func InitLogs() {

	lc, err := parseConfig("conf/log.yml")
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	if logger, err := New(&lc.Logger); err == nil {
		DefLogger = logger
	}
	ac, err := parseConfig("conf/access.yml")
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	if logger, err := New(&ac.Logger); err == nil {
		AccessLogger = logger
	}

}
func parseConfig(path string) (*LogConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &LogConfig{}
	if err := yaml.Unmarshal(data, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
