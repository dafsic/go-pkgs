package jwt

// Options .
type Options struct {
	// Intercept 拦截器
	*Intercept
}

type Option func(o *Options)

// defaultOptions .
func defaultOptions() *Options {
	return &Options{
		Intercept: &Intercept{enable: false},
	}
}

// Interceptor .
func Interceptor(ir *Intercept) Option {
	return func(o *Options) {
		ir.enable = true
		o.Intercept = ir
	}
}
