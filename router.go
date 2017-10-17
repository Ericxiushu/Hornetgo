package Hornetgo

import (
	"reflect"

	"net/http"

	"github.com/gorilla/mux"
)

type AppRouter struct {
	*mux.Router
}

func NewAppRouter() *AppRouter {
	return &AppRouter{
		mux.NewRouter(),
	}
}

//SetRoute SetRoute
func (c *AppRouter) SetRoute(path string, obj interface{}, action string) *AppRouter {

	r := &TempRouter{
		Path:   path,
		Obj:    obj,
		Action: action,
	}

	c.RegisterRouter(r)
	return c

}

func (c *AppRouter) RegisterRouter(r *TempRouter) *AppRouter {

	t := c.Handle(r.Path, r)
	if len(r.Methods) > 0 {
		t.Methods(r.Methods...)
	}
	return c
}

func (t *TempRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	hornetContent := &HornetContent{
		ResponseWriter: rw,
		Request:        r,
	}

	defer recoverPanic(hornetContent)

	reflectVal := reflect.ValueOf(t.Obj)
	tc := reflect.Indirect(reflectVal).Type()

	AppDebug("reflect obj :", tc)

	vc := reflect.New(tc)
	execController, ok := vc.Interface().(ControllerIntface)

	if !ok {
		AppDebug("not ok")
		panic("controller is not ControllerInterface")
	}

	execController.Init(hornetContent)

	execController.Start()

	var in []reflect.Value
	action := vc.MethodByName(t.Action)

	if action.IsValid() {
		action.Call(in)
	} else {
		panic("action is not exist")
	}
}

// recoverPanic recoverPanic
func recoverPanic(ctx *HornetContent) {

	if err := recover(); err != nil {
		AppDebug("Catch Panic : ", err)
		// render("error/errPage.html", map[interface{}]interface{}{"err_msg": err}, ctx)

		ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
		ctx.Write([]byte("error"))
	}

}
