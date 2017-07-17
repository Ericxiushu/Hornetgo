package Hornetgo

import (
	"encoding/json"
	"strings"

	"github.com/kataras/go-sessions"
	"github.com/pkg/errors"

	"github.com/qiangxue/fasthttp-routing"
)

type Contorller struct {
	Ctx       *routing.Context
	Data      map[interface{}]interface{}
	Session   sessions.Session
	YarWriter *YarWriter
}

type ControllerIntface interface {
	Init(*routing.Context, sessions.Session)
	Start()
}

func init() {
}

func (c *Contorller) Init(ctx *routing.Context, session sessions.Session) {

	c.Ctx = ctx
	c.YarWriter = &YarWriter{
		Ctx: ctx,
	}

	if strings.ToLower(string(c.Ctx.Request.Header.Peek("Content-Encoding"))) == "gzip" {
		body, err := c.Ctx.Request.BodyGunzip()
		if err == nil {
			c.Ctx.Request.SetBody(body)
		}
	}

	c.Data = make(map[interface{}]interface{}, 0)

	c.Session = session

}

func (c *Contorller) Start() {
	panic("Method Not Allowed")
}

// errorPage errorPage
func (c *Contorller) errorPage(str string) {

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

	doGzip(c.Ctx.RequestCtx)

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

func (c *Contorller) sendError(errCode int, errMsg string) {
	c.Data["json"] = map[string]interface{}{
		"err_code": errCode,
		"err_msg":  errMsg,
	}
	c.ServeJSON()
	return
}

func (c *Contorller) sendResult(result interface{}) {
	c.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "success",
		"data":     result,
	}
	return
}

func (c *Contorller) sendSuccess() {
	c.Data["json"] = map[string]interface{}{
		"err_code": 0,
		"err_msg":  "success",
		"data":     "",
	}
	c.ServeJSON()
	return
}

func (c *Contorller) sendOriInfo(result interface{}) {
	c.Data["json"] = result
	c.ServeJSON()
	return
}
