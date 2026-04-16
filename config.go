package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type netConfig struct {
	Ip   string
	Port string
	SRV  *http.Server
}

func (nc *netConfig) getURL() string {
	return "http://" + CFG.net.SRV.Addr
}

type osConfig struct {
	Uploads string
	TMPL    *template.Template
}

type Config struct {
	net netConfig
	os  osConfig
}

var (
	CFG Config
)

func init() {
	defaultUploads := filepath.Join(os.TempDir(), "uploads")
	flag.StringVar(&CFG.net.Port, "p", "8080", "port for serve files")
	flag.StringVar(&CFG.os.Uploads, "u", defaultUploads, "folder to save files")
	flag.Parse()
	CFG.net.SRV = &http.Server{
		Addr: GetLocalIP() + ":" + CFG.net.Port,
	}
	os.MkdirAll(CFG.os.Uploads, 0755)
}
