package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/e0m-ru/fileserver/internal/config"
	"github.com/e0m-ru/fileserver/internal/template"
	"github.com/e0m-ru/fileserver/internal/util"
)

// StartServer runs the file share server
func StartServer() error {
	var err error

	err = template.CollectTemplates()
	if err != nil {
		return err
	}

	config.Config.Net.SRV.Handler, err = CollectHandlers()
	if err != nil {
		return err
	}

	util.GenerateQRCode()

	var done = make(chan any, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	go func(srv_addr string, d chan any, s chan os.Signal) {
		sig := <-s
		log.Printf("\r\rreceived signal: %v\n", sig)
		config.Config.Net.SRV.Close()
		close(done)
	}(config.Config.Net.SRV.Addr, done, sig)

	fmt.Printf("upload dir: %v\n", config.Config.Os.Uploads)
	fmt.Printf("the fileserver is running on http://%v\n", config.Config.Net.SRV.Addr)
	if err := util.OpenURL("http://" + config.Config.Net.SRV.Addr); err != nil {
		return fmt.Errorf("❌ Error: %w\n", err)
	}
	if err := config.Config.Net.SRV.ListenAndServe(); err != nil {
		return fmt.Errorf("❌ Error: %w\n", err)
	}

	<-done

	return nil
}
