package main

import (
	"log"

	"github.com/EgorKo25/GophKeeper/internal/client"
	"github.com/EgorKo25/GophKeeper/internal/config"
	"github.com/EgorKo25/GophKeeper/internal/dialog"
	"github.com/EgorKo25/GophKeeper/pkg/mycrypto"
)

func main() {
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

	log.Println(dial.Run())

}
