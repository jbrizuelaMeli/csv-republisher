package customcontext

import "context"

func WithoutCancel(ctx context.Context) context.Context {
	return &uncancellableContext{
		Context: context.Background(),
		parent:  ctx,
	}
}

type uncancellableContext struct {
	context.Context
	parent context.Context
}

func (u *uncancellableContext) Value(key interface{}) interface{} {
	v := u.Context.Value(key)
	if v != nil {
		return v
	}
	return u.parent.Value(key)
}
