package Hornetgo

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/astaxie/beego/cache"
)

var (
	HornetInfo    *Hornet
	localCache, _ = cache.NewCache("memory", `{"interval":60}`)
)

func init() {
	HornetInfo = &Hornet{
		AppConfig: &Config{
			AppName:              "Hornet",
			RunMode:              RunModeDev,
			EnableGzip:           true,
			EnableSession:        true,
			EnableShowErrorsLine: true,
			Port:                 8080,
			WebConfig: WebConfig{
				ViewsPath: "/views",
			},
		},
		Router: NewRouter(),
	}
}

const (
	RunModeDev  = "dev"
	RunModeProd = "prod"
)

type Hornet struct {
	AppConfig *Config
	Router    *Router
}

func Run() error {

	Info("ListenAndServe port :", HornetInfo.AppConfig.Port)

	return fasthttp.ListenAndServe(fmt.Sprintf(":%d", HornetInfo.AppConfig.Port), HornetInfo.Router.HandleRequest)
}
