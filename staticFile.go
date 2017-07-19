package Hornetgo

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

func staticHander(ctx *fasthttp.RequestCtx) error {

	filePath, fileInfo, err := searchFile(ctx)

	if err != nil {
		return err
	}

	if filePath == "" || fileInfo == nil {
		return errors.New("empty")
	}

	if fileInfo.IsDir() {
		return errors.New("not a file")
	}

	fasthttp.ServeFile(ctx, filePath)

	return nil
}

// searchFile search the file by url path
// if none the static file prefix matches ,return notStaticRequestErr
func searchFile(ctx *fasthttp.RequestCtx) (string, os.FileInfo, error) {
	requestPath := filepath.ToSlash(filepath.Clean(string(ctx.Request.URI().Path())))
	// special processing : favicon.ico/robots.txt  can be in any static dir
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
		file := path.Join(".", requestPath)
		if fi, _ := os.Stat(file); fi != nil {
			return file, fi, nil
		}
		for _, staticDir := range HornetInfo.AppConfig.WebConfig.StaticDir {
			filePath := path.Join(staticDir, requestPath)
			if fi, _ := os.Stat(filePath); fi != nil {
				return filePath, fi, nil
			}
		}
		return "", nil, errors.New("403")
	}

	for prefix, staticDir := range HornetInfo.AppConfig.WebConfig.StaticDir {

		Error(requestPath, prefix, staticDir)

		if !strings.Contains(requestPath, prefix) {
			continue
		}
		if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
			continue
		}
		filePath := path.Join(staticDir, requestPath[len(prefix):])

		Error(filePath)

		if fi, err := os.Stat(filePath); fi != nil {
			return filePath, fi, err
		}
	}
	return "", nil, errors.New("403")
}
