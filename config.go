package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
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

	networks := GetAllLocalIPs()
	if len(networks) == 0 {
		fmt.Println("⚠️ No local networks found")
		return
	}

	var selected string
	prompt := &survey.Select{
		Message:  "Select network:",
		Options:  networks,
		PageSize: 10, // сколько показывать за раз
	}
	err := survey.AskOne(prompt, &selected)
	if err != nil {
		fmt.Printf("Input error: %v\n", err)
		return
	}

	fmt.Printf("✅ Selected: %s\n", selected)

	CFG.net.SRV = &http.Server{
		Addr: selected + ":" + CFG.net.Port,
	}

	os.MkdirAll(CFG.os.Uploads, 0755)
}
