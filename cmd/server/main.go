// package main is а package that contains server logic
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
	"net/http"
	"time"

	"github.com/EgorKo25/GophKeeper/pkg/auth"

	"github.com/EgorKo25/GophKeeper/internal/config"
	"github.com/EgorKo25/GophKeeper/internal/database"
	"github.com/EgorKo25/GophKeeper/internal/server/handlers"
	"github.com/EgorKo25/GophKeeper/internal/server/mymiddleware"
	"github.com/EgorKo25/GophKeeper/internal/server/myrouter"
)

var (
	buildVersion = "1.0.0"
	buildDate    = time.Now()
	buildCommit  = "Beta"
)

func main() {

	log.Printf("Версия приложения: %s\nДата сборки: %s\nТип версии: %s ", buildVersion, buildDate, buildCommit)

	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("config create error: %s", err)
	}

	db, err := database.NewManagerDB(cfg.DB)
	if err != nil {
		log.Fatalf("database constructor error: %s", err)
	}

	authentication := auth.NewAuth(cfg.RefreshToken)

	middle := mymiddleware.NewMyMiddleware(authentication, db)

	handler := handlers.NewHandler(db, authentication)

	router := myrouter.NewRouter(handler, middle)

	log.Println(http.ListenAndServe(cfg.Addr, router))
}
