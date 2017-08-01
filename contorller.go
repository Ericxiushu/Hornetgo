package Hornetgo

import (
	"encoding/json"

	"github.com/kataras/go-sessions"
	"github.com/pkg/errors"

	"github.com/qiangxue/fasthttp-routing"
)

type Contorller struct {
	Ctx     *routing.Context
	Data    map[interface{}]interface{}
	Session sessions.Session
}

type ControllerIntface interface {
	Init(*routing.Context, sessions.Session)
	Start()
}

func init() {
}

func (c *Contorller) Init(ctx *routing.Context, session sessions.Session) {

	c.Ctx = ctx

	c.Data = make(map[interface{}]interface{}, 0)

	c.Session = session

}

func (c *Contorller) Start() {
	panic("Method Not Allowed")
}

// errorPage errorPage
func (c *Contorller) ErrorPage(str string) {

	panic(str)

}

// Render Render
func (c *Contorller) Render(name string) {

	render(name, c.Data, c.Ctx.RequestCtx)

	return
}

// ServeJSON ServeJSON
func (c *Contorller) ServeJSON() {

	data, err := json.Marshal(c.Data["json"])
	if err == nil {
		c.Ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	} else {
		Error(errors.WithMessage(err, "json.Marshal"))
		data = []byte(err.Error())
	}

	c.Ctx.Write(data)

	return
}

// GetSession GetSession
func (c *Contorller) GetSession(key string) interface{} {
	return c.Session.Get(key)
}

// SetSession SetSession
func (c *Contorller) SetSession(key string, value interface{}) {
	c.Session.Set(key, value)
}

func (c *Contorller) SendError(errCode int, errMsg string) {
	c.Data["json"] = map[string]interface{}{
		"err_code": errCode,
		"err_msg":  errMsg,
	}
	c.ServeJSON()
	return
}

func (c *Contorller) SendResult(result interface{}) {
	c.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "success",
		"data":     result,
	}
	return
}

func (c *Contorller) SendSuccess() {
	c.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "success",
		"data":     "",
	}
	c.ServeJSON()
	return
}

func (c *Contorller) SendOriInfo(result interface{}) {
	c.Data["json"] = result
	c.ServeJSON()
	return
}

// Redirect Redirect
func (c *Contorller) Redirect(url string, code int) {
	c.Ctx.Redirect(url, code)
}

// GetString GetString
func (c *Contorller) GetString(key string) string {
	return string(c.Ctx.QueryArgs().Peek(key))
}
