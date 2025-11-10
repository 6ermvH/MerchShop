package logx

import "context"

type Nop struct{}

func (Nop) With(args ...any) Logger               { return Nop{} } //nolint:ireturn
func (Nop) Debug(context.Context, string, ...any) {}
func (Nop) Info(context.Context, string, ...any)  {}
func (Nop) Warn(context.Context, string, ...any)  {}
func (Nop) Error(context.Context, string, ...any) {}
