package Hornetgo

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Contorller struct {
	Ctx  *HornetContent
	Data map[interface{}]interface{}
}

type ControllerIntface interface {
	Init(*HornetContent)
	Start()
}

func init() {
}

func (c *Contorller) Init(ctx *HornetContent) {

	c.Ctx = ctx

	c.Data = make(map[interface{}]interface{}, 0)

}

func (c *Contorller) Start() {
	// panic("Method Not Allowed")
}

// errorPage errorPage
func (c *Contorller) ErrorPage(str string) {

	panic(str)

}

// Render Render
func (c *Contorller) Render(name string) {

	render(name, c.Data, c.Ctx)

	return
}

// ServeJSON ServeJSON
func (c *Contorller) ServeJSON() {

	data, err := json.Marshal(c.Data["json"])
	if err == nil {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	} else {
		Error(errors.WithMessage(err, "json.Marshal"))
		data = []byte(err.Error())
	}

	c.Ctx.Write(data)

	return
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
	c.ServeJSON()
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

// GetString GetString
func (c *Contorller) GetString(key string) string {

	result := c.Ctx.muxQuery[key]
	if len(result) < 1 {
		result = c.Ctx.URL.Query().Get(key)
	}

	return result
}
