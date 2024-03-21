package template

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"io/fs"
	"network-it/go-mail-autoconf/app/config"
	"text/template"
)

//go:embed templates
var templateFs embed.FS

type TemplateData struct {
	Env           *config.Config
	Domain        string
	DisplayName   string
	Address       string
	EmailProvider string
	Uuid          string
	Username      string
}

var TemplateMap = make(map[string]*template.Template)

func SetTemplate(fs fs.FS, name string, key string) {
	tpl, err := template.New(name).ParseFS(fs, name)
	if err != nil {
		log.Fatal(err)
	}
	TemplateMap[key] = tpl
}

func New() {
	tpls, _ := fs.Sub(fs.FS(templateFs), "templates")
	SetTemplate(tpls, "autodiscover.xml", "outlook")
	SetTemplate(tpls, "config-v1.1.xml", "thunderbird")
	SetTemplate(tpls, "email.mobileconfig.xml", "ios")
}

func NewData(address string) TemplateData {
	var cfg = config.NewConfig()
	return TemplateData{Env: cfg, Address: address, Username: address, DisplayName: address}
}

func Execute(key string, data TemplateData) (string, error) {
	var cfg = config.NewConfig()
	data.Env = cfg
	var tmplBytes bytes.Buffer
	tplt, ok := TemplateMap[key]
	if !ok {
		return "", fmt.Errorf("template %s not found", key)
	}
	err := tplt.Execute(&tmplBytes, data)
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
