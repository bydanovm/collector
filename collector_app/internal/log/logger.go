package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/mbydanov/collector/collector_app/internal/db/pgsql"
)

type LogEntry struct {
	Level      string    `gorm:"column:level"`
	Message    string    `gorm:"column:message"`
	Time       time.Time `gorm:"column:time"`
	Attributes string    `gorm:"column:attributes"`
}

type LoggerImpl interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type DBHandler struct {
	dbPgSQL pgsql.PgSQLInterface
}

func NewDBHandler(dbPgSQL pgsql.PgSQLInterface) *DBHandler {
	err := dbPgSQL.Db().AutoMigrate(LogEntry{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	return &DBHandler{dbPgSQL: dbPgSQL}
}

func (h *DBHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *DBHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	entry := &LogEntry{
		Level:      r.Level.String(),
		Message:    r.Message,
		Time:       r.Time,
		Attributes: fmt.Sprintf("%+v", attrs),
	}

	if err := h.dbPgSQL.Db().Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (h *DBHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *DBHandler) WithGroup(name string) slog.Handler {
	return h
}
