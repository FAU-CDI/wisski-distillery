package sql

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// logger implements "gorm/logger.Interface"
type gormLogger struct {
	Level                     logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func newGormLogger() logger.Interface {
	return &gormLogger{
		Level:                     logger.Info,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
}

// logLevelMap is the default
var logLevelMap = map[logger.LogLevel]zerolog.Level{
	logger.Silent:   zerolog.Disabled,
	logger.Error:    zerolog.ErrorLevel,
	logger.Warn:     zerolog.WarnLevel,
	logger.Info:     zerolog.DebugLevel,
	logger.Info + 1: zerolog.TraceLevel,
}

func (*gormLogger) NewEvent(ctx context.Context, level logger.LogLevel) *zerolog.Event {
	zl, ok := logLevelMap[level]
	if !ok {
		zl = zerolog.NoLevel
	}
	return zerolog.Ctx(ctx).WithLevel(zl)
}

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	new := *gl
	new.Level = level
	return &new
}

func (gl *gormLogger) Info(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Info {
		return
	}
	gl.NewEvent(ctx, logger.Info).Msgf(format, v...)
}
func (gl *gormLogger) Warn(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Warn {
		return
	}
	gl.NewEvent(ctx, logger.Warn).Msgf(format, v...)
}
func (gl *gormLogger) Error(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Error {
		return
	}
	gl.NewEvent(ctx, logger.Error).Msgf(format, v...)
}
func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.Level < logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && gl.Level >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !gl.IgnoreRecordNotFoundError):
		sql, rows := fc()
		src := utils.FileWithLineNum()

		gl.NewEvent(ctx, logger.Error).Str("src", src).Int64("rows", rows).Dur("elapsed", elapsed).Str("sql", sql).Msg("GORM")
	case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.Level >= logger.Warn:
		sql, rows := fc()
		src := utils.FileWithLineNum()
		gl.NewEvent(ctx, logger.Warn).Str("src", src).Int64("rows", rows).Dur("elapsed", elapsed).Str("sql", sql).Msgf("GORM: Slow SQL >= %d", gl.SlowThreshold)
	case gl.Level == logger.Info:
		sql, rows := fc()
		src := utils.FileWithLineNum()

		gl.NewEvent(ctx, logger.Info+1).Str("src", src).Int64("rows", rows).Dur("elapsed", elapsed).Str("sql", sql).Msg("GORM")
	}
}
