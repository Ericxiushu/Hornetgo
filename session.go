package Hornetgo

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/kataras/go-sessions"
	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"
)

var mySessions sessions.Sessions

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

	db := redis.New(service.Config{Network: service.DefaultRedisNetwork,
		Addr:          fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:      "",
		Database:      "",
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   service.DefaultRedisIdleTimeout,
		Prefix:        "",
		MaxAgeSeconds: 3600})

	mySessionsConfig := sessions.Config{Cookie: HornetInfo.AppConfig.AppName,
		Expires:                     time.Duration(1) * time.Hour,
		DisableSubdomainPersistence: false,
	}
	mySessions = sessions.New(mySessionsConfig)

	mySessions.UseDatabase(db)

}
