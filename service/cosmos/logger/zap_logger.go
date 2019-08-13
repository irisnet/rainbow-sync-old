package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type (
	Field = zap.Field
)

var (
	zapLogger *zap.Logger

	// zap method
	Binary     = zap.Binary
	Bool       = zap.Bool
	Complex128 = zap.Complex128
	Complex64  = zap.Complex64
	Float64    = zap.Float64
	Float32    = zap.Float32
	Int        = zap.Int
	Int64      = zap.Int64
	Int32      = zap.Int32
	Int16      = zap.Int16
	Int8       = zap.Int8
	String     = zap.String
	Uint       = zap.Uint
	Uint64     = zap.Uint64
	Uint32     = zap.Uint32
	Uint16     = zap.Uint16
	Uint8      = zap.Uint8
	Time       = zap.Time
	Any        = zap.Any
	Duration   = zap.Duration
)

func Debug(msg string, fields ...Field) {
	defer sync()
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defer sync()
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defer sync()
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	defer sync()
	zapLogger.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	defer sync()
	zapLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	defer sync()
	zapLogger.Fatal(msg, fields...)
}

func With(fields ...Field) {
	defer sync()
	zapLogger.With(fields...)
}

func sync() {
	zapLogger.Sync()
}

func init() {
	var core zapcore.Core
	hook := lumberjack.Logger{
		Filename:   "./logs/sync.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     7,    //days
		Compress:   true, // disabled by default
		LocalTime:  true,
	}

	fileWriter := zapcore.AddSync(&hook)
	consoleDebugging := zapcore.Lock(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
		zapcore.NewCore(encoder, fileWriter, highPriority),
	)
	caller := zap.AddCaller()
	callerSkipOpt := zap.AddCallerSkip(1)
	// From a zapcore.Core, it's easy to construct a Logger.
	zapLogger = zap.New(core, caller, callerSkipOpt, zap.AddStacktrace(zap.ErrorLevel))
}
