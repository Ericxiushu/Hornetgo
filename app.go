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
		Admin: Admin{
			PprofCPU: false,
			PprofMem: false,
		},
	}

	HornetInfo = &Hornet{
		AppConfig: AppConfig,
		AppRouter: NewAppRouter(),
	}
}

func Run() error {

	checkBeforeRun()

	for _, v := range TempRouterList {
		HornetInfo.AppRouter.SetRoute(v.Path, v.Obj, v.Methods...)
	}

	hander := HornetInfo.AppRouter.HandleRequest
	if HornetInfo.AppConfig.EnableGzip {
		hander = fasthttp.CompressHandler(hander)
	}

	Info("ListenAndServe port :", HornetInfo.AppConfig.Port)

	return fasthttp.ListenAndServe(fmt.Sprintf(":%d", HornetInfo.AppConfig.Port), hander)
}

func checkBeforeRun() {

	SetSession()

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

// func admin() {

// 	var err error

// 	if HornetInfo.AppConfig.Admin.PprofCPU {
// 		HornetInfo.AppConfig.Admin.CPUFile, err = os.OpenFile("./cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
// 		if err != nil {
// 			log.Fatal(err)
// 		} else {
// 			pprof.StartCPUProfile(HornetInfo.AppConfig.Admin.CPUFile)
// 		}
// 	}

// 	if HornetInfo.AppConfig.Admin.PprofMem {
// 		HornetInfo.AppConfig.Admin.MemFile, err = os.OpenFile("./mem.out", os.O_RDWR|os.O_CREATE, 0644)
// 		if err != nil {
// 			log.Fatal(err)
// 		} else {
// 			pprof.WriteHeapProfile(HornetInfo.AppConfig.Admin.MemFile)
// 		}
// 	}

// }

// func closeAdmin() {
// 	pprof.StopCPUProfile()
// 	HornetInfo.AppConfig.Admin.MemFile.Close()
// 	HornetInfo.AppConfig.Admin.CPUFile.Close()
// }
