package Hornetgo

type Config struct {
	AppName              string
	runMode              RunMode
	Port                 int
	EnableGzip           bool
	EnableSession        bool
	EnableShowErrorsLine bool
	WebConfig            WebConfig
}

type RunMode string

const (
	RunModeDev  RunMode = "dev"
	RunModeProd RunMode = "prod"
)

func (c RunMode) ToString() string {
	return string(c)
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
