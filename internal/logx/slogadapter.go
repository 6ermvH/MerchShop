package logx

import (
	"context"
	"log/slog"
)

type Slog struct{ l *slog.Logger }

func NewSlog(l *slog.Logger) *Slog { return &Slog{l: l} }

func (s *Slog) With(args ...any) Logger { return &Slog{l: s.l.With(args...)} }
func (s *Slog) Debug(ctx context.Context, msg string, args ...any) {
	s.l.DebugContext(ctx, msg, args...)
}

func (s *Slog) Info(
	ctx context.Context,
	msg string,
	args ...any,
) {
	s.l.InfoContext(ctx, msg, args...)
}

func (s *Slog) Warn(
	ctx context.Context,
	msg string,
	args ...any,
) {
	s.l.WarnContext(ctx, msg, args...)
}
func (s *Slog) Error(ctx context.Context, msg string, args ...any) {
	s.l.ErrorContext(ctx, msg, args...)
}
