package Hornetgo

import (
	"reflect"

	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Router struct {
	*routing.Router
}

func NewRouter() *Router {
	return &Router{
		routing.New(),
	}
}

//SetRoute SetRoute
func (c *Router) SetRoute(path string, obj interface{}, methods ...string) *routing.Route {

	if len(methods) > 0 {
		switch strings.ToLower(methods[0]) {
		case "get":
			return c.Get(path, RegisterRouter(obj))
		case "post":
			return c.Post(path, RegisterRouter(obj))
		default:
			return c.Any(path, RegisterRouter(obj))
		}
	}

	return c.Any(path, RegisterRouter(obj))

}

func RegisterRouter(obj interface{}) func(ctx *routing.Context) error {

	return func(ctx *routing.Context) error {

		defer recoverPanic(ctx.RequestCtx)

		reflectVal := reflect.ValueOf(obj)
		t := reflect.Indirect(reflectVal).Type()

		Error("reflect obj :", t)

		vc := reflect.New(t)
		execController, ok := vc.Interface().(ControllerIntface)

		if !ok {
			Error("not ok")
			panic("controller is not ControllerInterface")
		}

		execController.Init(ctx)

		execController.Start()

		return nil
	}

}

// recoverPanic recoverPanic
func recoverPanic(ctx *fasthttp.RequestCtx) {

	if err := recover(); err != nil {
		Error("Catch Panic : ", err)
		ShowPage("error/errPage.html", map[interface{}]interface{}{"err_msg": err}, ctx)
	}

}
