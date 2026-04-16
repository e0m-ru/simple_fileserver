package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	// "github.com/skip2/go-qrcode"
)

//go:embed static
var staticFiles embed.FS

func main() {

	var err error
	err = collectTemplates(&CFG)
	if err != nil {
		log.Print(err)
	}

	CFG.net.SRV.Handler, err = collectHandlers()
	if err != nil {
		log.Print(err)
	}

	qrCode()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	var done = make(chan any, 1)
	go func(s chan os.Signal, srv *http.Server, d chan any) {
		<-s
		fmt.Print("\r")
		srv.Close()
		close(done)
	}(sig, CFG.net.SRV, done)

	fmt.Print("upload dir: ", CFG.os.Uploads+"\n")
	fmt.Print("the fileserver is running on http://", CFG.net.SRV.Addr+"\n")
	if err := openURL("http://" + CFG.net.SRV.Addr); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
	if err := CFG.net.SRV.ListenAndServe(); err != nil {
		log.Print(err)
	}

	<-done
}
