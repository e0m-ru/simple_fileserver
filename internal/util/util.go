package util

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/e0m-ru/fileserver/internal/config"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func GenerateQRCode() {
	qrc, err := qrcode.New(config.Config.Net.GetURL())
	if err != nil {
		log.Print(err)
	}

	w, err := standard.New(config.Config.Os.Uploads + "/qrcode.jpeg")
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}

	// save file
	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}

func OpenURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}
