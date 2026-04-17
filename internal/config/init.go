package config

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

func init() {
	// Skip initialization during tests
	if os.Getenv("TESTING") == "1" {
		// Set defaults for testing
		Config.Net.Port = "8080"
		Config.Os.Uploads = filepath.Join(os.TempDir(), "uploads")
		return
	}

	defaultUploads := filepath.Join(os.TempDir(), "uploads")

	flag.StringVar(&Config.Net.Port, "p", "8080", "port for serve files")
	flag.StringVar(&Config.Os.Uploads, "u", defaultUploads, "folder to save files")
	flag.Parse()

	networks, err := GetAllLocalIPs()
	if err != nil {
		log.Print(err)
	}

	if len(networks) == 0 {
		err := fmt.Errorf("⚠️ No local networks found: %v", networks)
		log.Fatal(err)
	}

	var selected string
	prompt := &survey.Select{
		Message:  "Select network:",
		Options:  networks,
		PageSize: 10,
	}
	err = survey.AskOne(prompt, &selected)
	if err != nil {
		fmt.Printf("Input error: %v\n", err)
		return
	}

	fmt.Printf("✅ Selected: %s\n", selected)

	Config.Net.SRV = &http.Server{
		Addr: selected + ":" + Config.Net.Port,
	}

	os.MkdirAll(Config.Os.Uploads, 0755)
}
