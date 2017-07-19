package Hornetgo

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"strings"

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
				StaticDir: map[string]string{"/static": "static"},
			},
		},
		AppRouter: NewAppRouter(),
	}
}

type Hornet struct {
	AppConfig *Config
	AppRouter *AppRouter
}

func Run() error {

	checkBeforeRun()

	Info("ListenAndServe port :", HornetInfo.AppConfig.Port)

	hander := HornetInfo.AppRouter.HandleRequest
	if HornetInfo.AppConfig.EnableGzip {
		hander = fasthttp.CompressHandler(hander)
	}

	return fasthttp.ListenAndServe(fmt.Sprintf(":%d", HornetInfo.AppConfig.Port), hander)
}

func checkBeforeRun() {

	// 检测session
	if HornetInfo.AppConfig.EnableSession && mySessions == nil {
		panic("manager session error")
	}

	// 注册静态资源路由
	for path := range HornetInfo.AppConfig.WebConfig.StaticDir {
		path = strings.TrimSuffix(path, "/") + "/*"
		HornetInfo.AppRouter.Any(path, serverStaticRouter)
	}

}

func Router(path string, obj interface{}, methods ...string) {
	HornetInfo.AppRouter.SetRoute(path, obj, methods...)
}
