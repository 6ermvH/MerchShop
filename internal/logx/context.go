package logx

import "context"

type ctxKey struct{}

func IntoContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) Logger {
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(Logger); ok && l != nil {
			return l
		}
	}
	return Nop{}
}

type Nop struct{}

func (Nop) With(args ...any) Logger               { return Nop{} }
func (Nop) Debug(context.Context, string, ...any) {}
func (Nop) Info(context.Context, string, ...any)  {}
func (Nop) Warn(context.Context, string, ...any)  {}
func (Nop) Error(context.Context, string, ...any) {}
