package Hornetgo

import (
	"bufio"
	"bytes"

	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// YarWriter YarWriter
type YarWriter struct {
	Ctx *routing.Context
}

// Write Write
func (c *YarWriter) Write(p []byte) (int, error) {

	c.Ctx.Write(p)
	doGzip(c.Ctx.RequestCtx)

	return len(p), nil

}

// doGzip doGzip
func doGzip(ctx *fasthttp.RequestCtx) {

	if HornetInfo.AppConfig.EnableGzip {
		w := &bytes.Buffer{}
		bw := bufio.NewWriter(w)
		err := ctx.Response.WriteGzip(bw)

		if err != nil {
			Error(errors.WithMessage(err, "WriteGzip"))
		}

		bw.Flush()

	}
	return
}
