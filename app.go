package Hornetgo

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"strings"

	"github.com/astaxie/beego/cache"
)

var (
	HornetInfo     *Hornet
	AppConfig      *Config
	TempRouterList []TempRouter
	localCache, _  = cache.NewCache("memory", `{"interval":60}`)
)

type Hornet struct {
	AppConfig *Config
	AppRouter *AppRouter
}

func init() {

	AppConfig = &Config{
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
	}

}

func Run() error {

	HornetInfo = &Hornet{
		AppConfig: AppConfig,
		AppRouter: NewAppRouter(),
	}

	checkBeforeRun()

	for _, v := range TempRouterList {
		HornetInfo.AppRouter.SetRoute(v.Path, v.Obj, v.Methods...)
	}

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

	// 检测runmodel
	if HornetInfo.AppConfig.RunMode != RunModeDev && HornetInfo.AppConfig.RunMode != RunModeProd {
		panic(" not allowed runmodel ")
	}

}

func Router(path string, obj interface{}, methods ...string) {
	// HornetInfo.AppRouter.SetRoute(path, obj, methods...)

	item := TempRouter{
		Path:    path,
		Obj:     obj,
		Methods: methods,
	}

	TempRouterList = append(TempRouterList, item)
}
