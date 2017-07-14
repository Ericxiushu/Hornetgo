package Hornetgo

import (
	"reflect"

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

func (c *Router) SetAny(path string, obj interface{}) *routing.Route {
	return c.Any(path, StartRouter(obj))
}

func StartRouter(obj interface{}) func(ctx *routing.Context) error {

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
		Error("panic : ", err)
		ShowPage("error/errPage.html", map[interface{}]interface{}{"err_msg": err}, ctx)
	}

}
