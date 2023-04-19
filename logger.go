package logger

import (
	"context"
	"errors"
	"log"

	"github.com/TheZeroSlave/zapsentry"

	"pkg.moe/pkg/contexts"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger       *zap.Logger
	loggerConfig zap.Config
)

func Init(lvString string) {
	var (
		err error
		lvl zapcore.Level
	)

	if lvl, err = getLoggerLevel(lvString); err != nil {
		log.Fatalln("failed to initalize logger due to:", err)
	}

	if lvl == zapcore.DebugLevel {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.Level = zap.NewAtomicLevelAt(lvl)
	logger, err = loggerConfig.Build(
		zap.AddStacktrace(zapcore.PanicLevel),
	)

	if err != nil {
		log.Fatalln("failed to initialize logger due to:", err)
	}

}

func Get() *zap.SugaredLogger {
	if logger == nil {
		Init("error")
	}

	return logger.Sugar()
}

func GetWithContext(ctx context.Context) *zap.Logger {
	if logger == nil {
		Init("error")
	}

	var zapFields []zap.Field

	// basic info
	uid := contexts.Int64(ctx, "uid")
	if uid != 0 {
		zapFields = append(zapFields, zap.Int64("uid", uid))
	}

	issuer := contexts.String(ctx, "issuer")
	if issuer != "" {
		zapFields = append(zapFields, zap.String("issuer", issuer))
	}

	ip := contexts.String(ctx, "ip")
	if ip != "" {
		zapFields = append(zapFields, zap.String("ip", ip))
	}

	device_ip := contexts.String(ctx, "device_ip")
	if device_ip != "" {
		zapFields = append(zapFields, zap.String("device_ip", device_ip))
	}

	device_id := contexts.String(ctx, "device_id")
	if device_id != "" {
		zapFields = append(zapFields, zap.String("device_id", device_id))
	}

	app_version := contexts.String(ctx, "app_version")
	if app_version != "" {
		zapFields = append(zapFields, zap.String("app_version", app_version))
	}

	platform := contexts.String(ctx, "platform")
	if platform != "" {
		zapFields = append(zapFields, zap.String("platform", platform))
	}

	return logger.With(zapFields...)
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func FieldError(err error) zap.Field {
	return zap.Error(err)
}

func getLoggerLevel(lvString string) (zapcore.Level, error) {
	var (
		lvl zapcore.Level
	)

	if err := lvl.UnmarshalText([]byte(lvString)); err != nil {
		return lvl, err
	}

	return lvl, nil
}

func SetLevel(lvString string) error {
	lvl, err := getLoggerLevel(lvString)
	if err != nil {
		return errors.New("failed to set logger level due to:" + err.Error())
	}
	loggerConfig.Level = zap.NewAtomicLevelAt(lvl)
	newLogger, err := loggerConfig.Build(
		zap.AddStacktrace(zapcore.PanicLevel),
	)

	if err != nil {
		return errors.New("failed to set logger level due to:" + err.Error())
	}

	logger = newLogger
	return nil

}

func SetSentryLogger(DSN string) error {
	cfg := zapsentry.Configuration{
		Level: zapcore.ErrorLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(DSN))
	if err != nil {
		return err
	}
	logger = zapsentry.AttachCoreToLogger(core, logger)

	return nil
}
