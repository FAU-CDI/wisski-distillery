package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
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

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	new := *gl
	new.Level = level
	return &new
}

func (gl *gormLogger) Info(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Info {
		return
	}

	wdlog.Of(ctx).Info(
		"gorm",
		"message", fmt.Sprintf(format, v...),
	)
}
func (gl *gormLogger) Warn(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Warn {
		return
	}
	wdlog.Of(ctx).Warn(
		"gorm",
		"message", fmt.Sprintf(format, v...),
	)
}
func (gl *gormLogger) Error(ctx context.Context, format string, v ...interface{}) {
	if gl.Level < logger.Error {
		return
	}
	wdlog.Of(ctx).Error(
		"gorm",
		"message", fmt.Sprintf(format, v...),
	)
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

		wdlog.Of(ctx).Error(
			"gorm",

			"src", src,
			"rows", rows,
			"elapsed", elapsed,
			"sql", sql,
		)
	case elapsed > gl.SlowThreshold && gl.SlowThreshold != 0 && gl.Level >= logger.Warn:
		sql, rows := fc()
		src := utils.FileWithLineNum()
		wdlog.Of(ctx).Warn(
			fmt.Sprintf("gorm: Slow SQL >= %d", gl.SlowThreshold),

			"src", src,
			"rows", rows,
			"elapsed", elapsed,
			"sql", sql,
		)
	case gl.Level == logger.Info:
		sql, rows := fc()
		src := utils.FileWithLineNum()

		wdlog.Of(ctx).Debug(
			"gorm",

			"src", src,
			"rows", rows,
			"elapsed", elapsed,
			"sql", sql,
		)
	}
}
