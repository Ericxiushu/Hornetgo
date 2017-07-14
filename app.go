package Hornetgo

import "github.com/astaxie/beego/cache"

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
