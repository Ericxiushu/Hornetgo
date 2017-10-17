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

func (c *AppRouter) SetRouter(pathPre string, r []*TempRouter) *AppRouter {

	var t *mux.Router
	if len(pathPre) > 0 {
		t = c.PathPrefix(pathPre).Subrouter()

		for _, v := range r {
			newR := t.Handle(v.Path, v)
			if len(v.Methods) > 0 {
				newR.Methods(v.Methods...)
			}
		}

	} else {

		for _, v := range r {
			newR := c.Handle(v.Path, v)
			if len(v.Methods) > 0 {
				newR.Methods(v.Methods...)
			}
		}

	}

	return c
}

func (t *TempRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	hornetContent := &HornetContent{
		ResponseWriter: rw,
		Request:        r,
		muxQuery:       mux.Vars(r),
	}

	hornetContent.CopyBody()

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

		ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		ctx.Write([]byte("error"))
	}

}
