// package main is а package that contains client logic
//
// Build command:
//
//	go build main.go
//
// Run command:
//
//	go run main.go
package main

import (
	"log"
	"time"

	"github.com/EgorKo25/GophKeeper/internal/client"
	"github.com/EgorKo25/GophKeeper/internal/config"
	"github.com/EgorKo25/GophKeeper/internal/dialog"
	"github.com/EgorKo25/GophKeeper/pkg/mycrypto"
)

var (
	buildVersion = "1.0.0"
	buildDate    = time.Now()
	buildCommit  = "Beta"
)

func main() {

	log.Printf("Версия приложения: %s\nДата сборки: %s\nТип версии: %s ", buildVersion, buildDate, buildCommit)

	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
	}

	enc, err := mycrypto.NewCrypto(cfg.Secret)
	if err != nil {
		log.Fatal(err)
	}

	c := client.NewClient(cfg.AddrServ)

	dial := dialog.NewDialogManager(c, enc)

	dial.SayHello()

	log.Println(dial.Run())

}
