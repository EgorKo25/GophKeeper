package main

import (
	"log"
	"net/http"

	"github.com/EgorKo25/GophKeeper/pkg/mycrypto"

	"github.com/EgorKo25/GophKeeper/pkg/auth"

	"github.com/EgorKo25/GophKeeper/internal/config"
	"github.com/EgorKo25/GophKeeper/internal/database"
	"github.com/EgorKo25/GophKeeper/internal/server/handlers"
	"github.com/EgorKo25/GophKeeper/internal/server/mymiddleware"
	"github.com/EgorKo25/GophKeeper/internal/server/myrouter"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("config create error: %s", err)
	}

	db, err := database.NewManagerDB(cfg.DB)
	if err != nil {
		log.Fatalf("database constructor error: %s", err)
	}

	encrypt, err := mycrypto.NewCrypto()
	if err != nil {
		log.Fatalf("encrypt constructor error: %s", err)
	}

	authentication := auth.NewAuth(cfg.RefreshToken)

	middle := mymiddleware.NewMyMiddleware(authentication, db)

	handler := handlers.NewHandler(db, authentication, encrypt)

	router := myrouter.NewRouter(handler, middle)

	log.Println(http.ListenAndServe(cfg.Addr, router))
}
