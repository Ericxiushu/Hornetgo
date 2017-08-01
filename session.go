package Hornetgo

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/clevergo/sessions"
)

var (
	store sessions.Store
)

func init() {

	if !HornetInfo.AppConfig.EnableSession {
		return
	}

	type EnvConfig struct {
		Host string `env:"WXHOST_SYSTEM_REDIS_HOST"`
		Port int    `env:"WXHOST_SYSTEM_REDIS_PORT"`
	}

	cfg := EnvConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	if len(cfg.Host) < 1 || cfg.Port < 1 {
		panic("redis config error")
	}

	store, err = sessions.NewRediStore(10, "tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}

}
