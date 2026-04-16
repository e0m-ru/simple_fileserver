package main

import (
	"flag"
	"fmt"
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
	fmt.Printf("Select network:")
	networks := GetAllLocalIPs()
	for i, a := range networks {
		fmt.Println(i, ": ", a)
	}
	var n int
	fmt.Scan(&n)
	fmt.Printf("Network: %v\n", n)
	CFG.net.SRV = &http.Server{
		Addr: networks[n] + ":" + CFG.net.Port,
	}
	os.MkdirAll(CFG.os.Uploads, 0755)
}
