package waitgroup

import "context"

type (
	// Option defines a functional option type for the Wait method.
	Option func(*option)

	option struct {
		WithContext context.Context
	}

	options []Option
)

// WithContext will supply the given context to Wait for use with
// timeouts.
func WithContext(ctx context.Context) Option {
	return func(o *option) {
		o.WithContext = ctx
	}
}
