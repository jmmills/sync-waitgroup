package waitgroup

import "context"

type (
	Option func(*option)

	option struct {
		WithContext context.Context
	}

	options []Option
)

func WithContext(ctx context.Context) Option {
	return func(o *option) {
		o.WithContext = ctx
	}
}
