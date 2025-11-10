package logx

import "context"

type ctxKey struct{}

func IntoContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) Logger { //nolint:ireturn
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(Logger); ok && l != nil {
			return l
		}
	}

	return Nop{}
}
