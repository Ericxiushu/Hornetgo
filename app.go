package Hornetgo

import (
	"fmt"
	"strings"

	"net/http"

	"github.com/astaxie/beego/cache"
	"github.com/gorilla/mux"
)

var (
	HornetInfo    *Hornet
	AppConfig     *Config
	TempRouterMap = make(map[string][]*TempRouter, 0)
	localCache, _ = cache.NewCache("memory", `{"interval":60}`)
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

	for k, v := range TempRouterMap {
		HornetInfo.AppRouter.SetRouter(k, v)
	}

	fmt.Println(TempRouterMap)

	HornetInfo.AppRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			fmt.Println(err)

			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			fmt.Println(err)

			return err
		}
		fmt.Println(strings.Join(m, ","), t, p)
		return nil
	})

	Info("ListenAndServe port :", HornetInfo.AppConfig.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", HornetInfo.AppConfig.Port), HornetInfo.AppRouter)

}

func checkBeforeRun() {

	// 检测runmodel
	if HornetInfo.AppConfig.RunMode != RunModeDev && HornetInfo.AppConfig.RunMode != RunModeProd {
		panic(" not allowed runmodel ")
	}

}

func Router(pathPre, path string, obj interface{}, action string, methods ...string) {

	item := &TempRouter{
		Path:    path,
		Obj:     obj,
		Methods: methods,
		Action:  action,
	}

	if _, ok := TempRouterMap[pathPre]; !ok {
		TempRouterMap[pathPre] = make([]*TempRouter, 0)
	}

	TempRouterMap[pathPre] = append(TempRouterMap[pathPre], item)

}
