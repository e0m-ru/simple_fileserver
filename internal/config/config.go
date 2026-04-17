package config

import (
	"embed"
	"html/template"
	"net"
	"net/http"
)

type NetConfig struct {
	Ip   string
	Port string
	SRV  *http.Server
}

func (nc *NetConfig) GetURL() string {
	return "http://" + nc.SRV.Addr
}

type OsConfig struct {
	Uploads string
	TMPL    *template.Template
}

type AppConfig struct {
	Net *NetConfig
	Os  *OsConfig
}

var (
	Config = &AppConfig{
		Net: &NetConfig{},
		Os:  &OsConfig{},
	}
)

//go:embed static/css/*.css static/js/*.js
var StaticFiles embed.FS

// GetAllLocalIPs returns a list of the system's unicast interface addresses
func GetAllLocalIPs() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips, nil
}
