package Hornetgo

import (
	"reflect"

	"strings"

	"github.com/kataras/go-sessions"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type AppRouter struct {
	*routing.Router
}

func NewAppRouter() *AppRouter {
	return &AppRouter{
		routing.New(),
	}
}

func serverStaticRouter(ctx *routing.Context) error {
	return staticHander(ctx.RequestCtx)
}

//SetRoute SetRoute
func (c *AppRouter) SetRoute(path string, obj interface{}, methods ...string) *routing.Route {

	if len(methods) > 0 {

		for _, v := range methods {
			switch strings.ToLower(v) {
			case "get":
				return c.Get(path, RegisterRouter(obj))
			case "post":
				return c.Post(path, RegisterRouter(obj))
			case "any":
				return c.Any(path, RegisterRouter(obj))
			default:
				panic("not support ")
			}
		}

	}

	return c.Any(path, RegisterRouter(obj))

}

func RegisterRouter(obj interface{}) func(ctx *routing.Context) error {

	return func(ctx *routing.Context) error {

		defer recoverPanic(ctx.RequestCtx)

		reflectVal := reflect.ValueOf(obj)
		t := reflect.Indirect(reflectVal).Type()

		AppDebug("reflect obj :", t)

		vc := reflect.New(t)
		execController, ok := vc.Interface().(ControllerIntface)

		if !ok {
			AppDebug("not ok")
			panic("controller is not ControllerInterface")
		}

		var session sessions.Session
		if HornetInfo.AppConfig.EnableSession {
			session = mySessions.StartFasthttp(ctx.RequestCtx)

			defer func() {
				// todo : 关闭session
			}()
		}

		execController.Init(ctx, session)

		execController.Start()

		return nil
	}

}

// recoverPanic recoverPanic
func recoverPanic(ctx *fasthttp.RequestCtx) {

	if err := recover(); err != nil {
		AppDebug("Catch Panic : ", err)
		render("error/errPage.html", map[interface{}]interface{}{"err_msg": err}, ctx)
	}

}
