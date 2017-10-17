package Hornetgo

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	AppName              string
	RunMode              string
	Port                 int
	EnableGzip           bool
	EnableSession        bool
	EnableShowErrorsLine bool
	WebConfig            WebConfig
	Admin                Admin
}

const (
	RunModeDev  = "dev"
	RunModeProd = "prod"
)

type Admin struct {
	PprofMem bool
	PprofCPU bool
	CPUFile  *os.File
	MemFile  *os.File
}

// WebConfig holds web related config
type WebConfig struct {
	AutoRender             bool
	EnableDocs             bool
	FlashName              string
	FlashSeparator         string
	DirectoryIndex         bool
	StaticDir              map[string]string
	StaticExtensionsToGzip []string
	TemplateLeft           string
	TemplateRight          string
	ViewsPath              string
	EnableXSRF             bool
	XSRFKey                string
	XSRFExpire             int
}

type TempRouter struct {
	Path    string
	Obj     interface{}
	Action  string
	Methods []string
}

type HornetContent struct {
	http.ResponseWriter
	*http.Request
	Body     []byte
	muxQuery map[string]string
}

func (c *HornetContent) Write(b []byte) (int, error) {

	buf := &bytes.Buffer{}

	buf.Write(b)

	io.Copy(c.ResponseWriter, buf)

	return len(b), nil
}

func (c *HornetContent) CopyBody() {
	c.Body, _ = ioutil.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	return
}
