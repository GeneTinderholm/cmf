package config

import (
	"context"
	"runtime"
)

var defaultConfig Config

func init() {
	defaultConfig.Parallelism = runtime.GOMAXPROCS(-1)
}

type configKey string
var ck configKey = "configKey"

type Config struct {
	Parallelism int
}

func WithConfig(ctx context.Context, c Config) context.Context {
	return context.WithValue(ctx, ck, c)
}
func GetConfig(ctx context.Context) Config {
	if cfg, ok := ctx.Value(ck).(Config); ok {
		return cfg
	} else {
		return defaultConfig
	}
}
