package Hornetgo

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/kataras/go-sessions"
	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"
	"github.com/valyala/fasthttp"

	"github.com/qiangxue/fasthttp-routing"
)

var mySessionsConfig = sessions.Config{Cookie: "mysessioncookieid",
	Expires:                     time.Duration(1) * time.Hour,
	DisableSubdomainPersistence: false,
}

var mySessions = sessions.New(mySessionsConfig)

type Contorller struct {
	Ctx       *routing.Context
	Data      map[interface{}]interface{}
	Session   sessions.Session
	YarWriter *YarWriter
}

type ControllerIntface interface {
	Init(*routing.Context)
	Start()
}

func init() {
	startSession()
}

func (c *Contorller) Init(ctx *routing.Context) {

	c.Ctx = ctx
	c.YarWriter = &YarWriter{
		Ctx: ctx,
	}

	if strings.ToLower(string(c.Ctx.Request.Header.Peek("Content-Encoding"))) == "gzip" {
		body, err := c.Ctx.Request.BodyGunzip()
		if err == nil {
			c.Ctx.Request.SetBody(body)
		}
	}

	c.Data = make(map[interface{}]interface{}, 0)
	c.Session = mySessions.StartFasthttp(c.Ctx.RequestCtx)

}

func (c *Contorller) Start() {
	panic("Method Not Allowed")
}

// errorPage errorPage
func (c *Contorller) errorPage(str string) {

	panic(str)

}

func startSession() {

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

	mySessions.UseDatabase(db)

}

// Render Render
func (c *Contorller) Render(name string) {

	render(name, c.Data, c.Ctx.RequestCtx)

	return
}

func render(name string, data interface{}, ctx *fasthttp.RequestCtx) {

	body, err := GetPageByTemplate(name, data)
	if err != nil {
		Error(err)
		body = []byte("not find")
	}

	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Write(body)

	doGzip(ctx)

	return
}
