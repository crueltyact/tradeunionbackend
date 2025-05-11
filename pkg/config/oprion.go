package config

const (
	defaultConfigPath = "./config/config.json"
)

type cfgOptions struct {
	cfgPath         string
	disableDefaults bool
}

type Option interface {
	apply(*cfgOptions)
}

type cfgPathOption struct {
	cfgPath string
}

func (o *cfgPathOption) apply(opts *cfgOptions) {
	opts.cfgPath = o.cfgPath
}

func WithConfigPath(cfgPath string) Option {
	return &cfgPathOption{cfgPath: cfgPath}
}

type defaultValuesOpt struct {
	disableDefaults bool
}

func (o *defaultValuesOpt) apply(opts *cfgOptions) {
	opts.disableDefaults = o.disableDefaults
}

func newDefaultConfigOptions() *cfgOptions {
	return &cfgOptions{
		cfgPath:         defaultConfigPath,
		disableDefaults: false,
	}
}
