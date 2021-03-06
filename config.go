package Hornetgo

import (
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
	Methods []string
}
