package logger

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger is a run log instance
	Log *zap.Logger = zap.NewExample()
)

func init() {

	cfg := Config{
		EncodeLogsAsJson:   true,
		FileLoggingEnabled: true,
		Directory:          "/logs/", // 要斜杠
		Filename:           "run.log",
		MaxSize:            512,
		MaxBackups:         30,
		MaxAge:             7,
	}
	Log = NewLogger(cfg)

}

func Logger() *zap.Logger {
	return DefaultZapLogger
}

// zaplog.Trace(req.Id).Info("11111111111")
func Trace(id string) *zap.Logger {

	return DefaultZapLogger.With(zap.String("id", id))

}

func SetTraceId(id string) *zap.Logger {

	DefaultZapLogger = DefaultZapLogger.With(zap.String("id", id))

	return DefaultZapLogger

}

// includes the error message in the final log output.
func Reflect(key string, val interface{}) zapcore.Field {
	return zap.Reflect(key, val)
}

// String constructs a field with the given key and value.
func String(key string, val string) zapcore.Field {
	return zap.String(key, val)
}

// Err
func Err(err error) zapcore.Field {
	return zap.Error(err)
}

// Configuration for logging
type Config struct {
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int

	Debug bool
}

// How to log, by example:
// logger.Info("Importing new file, zap.String("source", filename), zap.Int("size", 1024))
// To log a stacktrace:
// logger.Error("It went wrong, zap.Stack())

// DefaultZapLogger is the default logger instance that should be used to log
// It's assigned a default value here for tests (which do not call log.Configure())
var DefaultZapLogger = newZapLogger(false, os.Stdout)

// Debug Log a message at the debug level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Debug(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Debug(msg, fields...)
}

// Info log a message at the info level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Info(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Info(msg, fields...)
}

// Warn log a message at the warn level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Warn(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Warn(msg, fields...)
}

// Error Log a message at the error level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Error(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Error(msg, fields...)
}

// Panic Log a message at the Panic level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Panic(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Panic(msg, fields...)
}

// Fatal Log a message at the fatal level. Messages include any context that's
// accumulated on the logger, as well as any fields added at the log site.
//
// Use zap.String(key, value), zap.Int(key, value) to log fields. These fields
// will be marshalled as JSON in the logfile and key value pairs in the console!
func Fatal(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.Fatal(msg, fields...)
}

// AtLevel logs the message at a specific log level
func AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) {
	switch level {
	case zapcore.DebugLevel:
		Debug(msg, fields...)
	case zapcore.PanicLevel:
		Panic(msg, fields...)
	case zapcore.ErrorLevel:
		Error(msg, fields...)
	case zapcore.WarnLevel:
		Warn(msg, fields...)
	case zapcore.InfoLevel:
		Info(msg, fields...)
	case zapcore.FatalLevel:
		Fatal(msg, fields...)
	default:
		Warn("Logging at unkown level", zap.Any("level", level))
		Warn(msg, fields...)
	}
}

// TID trace_id的key
const (
	//TID gin: key of trace_id
	TID = "trace_id"
	//KeyTraceID trpc: key of trace_id
	KeyTraceID = ""
)

type logFunc func(msg string, fields ...zapcore.Field)

// 装饰器,日志打印trace_id
func withTrace(c context.Context, log logFunc) logFunc {
	return func(msg string, fields ...zapcore.Field) {
		//gin context
		if v, ok := c.(*gin.Context); ok { //不能直接c==nil判断! c为interface,即使被nil的*gin.Context赋值,也不等于nil
			if v == nil { //指针需要判断空
				log(msg, fields...)
				return
			}
			if traceID, ok := c.Value(TID).(string); ok {
				//获取到trace_id,追加在日志最后
				add := zap.String(TID, traceID)
				fields = append(fields, add)
			}
		}
		if c == nil {
			log(msg, fields...)
			return
		}
		//rpc context
		if traceID, ok := c.Value(KeyTraceID).(string); ok {
			//获取到trace_id,追加在日志最后
			add := zap.String(TID, traceID)
			fields = append(fields, add)
		}

		//因为装饰器让调用层级+1了,所以这里要skip两层,才能打印正确代码调用位置
		log(msg, fields...)
	}
}

// DebugT pring debug log with trace_id
func DebugT(c context.Context, msg string, fields ...zapcore.Field) {
	f := withTrace(c, debug2)
	f(msg, fields...)
}

// InfoT pring info log with trace_id
func InfoT(c context.Context, msg string, fields ...zapcore.Field) {
	f := withTrace(c, info2)
	f(msg, fields...)
}

// WarnT pring warn log with trace_id
func WarnT(c context.Context, msg string, fields ...zapcore.Field) {
	f := withTrace(c, warn2)
	f(msg, fields...)
}

// ErrorT pring error log with trace_id
func ErrorT(c context.Context, msg string, fields ...zapcore.Field) {
	f := withTrace(c, error2)
	f(msg, fields...)
}

func debug2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Debug(msg, fields...)
}

func info2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Info(msg, fields...)
}

func warn2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Warn(msg, fields...)
}

func error2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Error(msg, fields...)
}

// Panic func
func panic2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Panic(msg, fields...)
}

// Fatal func
func fatal2(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(2)).Fatal(msg, fields...)
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/auth-service/auth-service.log and
// will be rolled when it reaches 20MB with a maximum of 1 backup.
func Configure(config Config) {
	writers := []zapcore.WriteSyncer{}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	DefaultZapLogger = newZapLogger(config.EncodeLogsAsJson, zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(DefaultZapLogger)

}

func newRollingFile(config Config) zapcore.WriteSyncer {
	if err := os.MkdirAll(config.Directory, 0775); err != nil {
		Error("failed create log directory", zap.Error(err), zap.String("path", config.Directory))
		return nil
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Directory + config.Filename,
		MaxSize:    config.MaxSize,    //megabytes
		MaxAge:     config.MaxAge,     //days
		MaxBackups: config.MaxBackups, //files
	})
}

func newZapLogger(encodeAsJSON bool, output zapcore.WriteSyncer) *zap.Logger {

	encCfg := zapcore.EncoderConfig{
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

	encoder := zapcore.NewConsoleEncoder(encCfg)
	if encodeAsJSON {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	return zap.New(zapcore.NewCore(encoder, output, zap.NewAtomicLevel()), zap.AddCaller(), zap.AddCallerSkip(1))
}

// 自定义 instance
func NewLogger(config Config) *zap.Logger {

	writers := []zapcore.WriteSyncer{}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	lg := newZapLogger(config.EncodeLogsAsJson, zapcore.NewMultiWriteSyncer(writers...))
	zap.RedirectStdLog(DefaultZapLogger)

	return lg

}
